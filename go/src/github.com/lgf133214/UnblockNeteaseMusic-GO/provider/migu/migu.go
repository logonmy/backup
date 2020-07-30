package migu

import (
	"bytes"
	"crypto/md5"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/lgf133214/UnblockNeteaseMusic-GO/common"
	"github.com/lgf133214/UnblockNeteaseMusic-GO/network"
	"github.com/lgf133214/UnblockNeteaseMusic-GO/processor/crypto"
	"github.com/lgf133214/UnblockNeteaseMusic-GO/utils"
	"math"
	"net/http"
	"net/url"
	"strings"
)

var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC8asrfSaoOb4je+DSmKdriQJKW
VJ2oDZrs3wi5W67m3LwTB9QVR+cE3XWU21Nx+YBxS0yun8wDcjgQvYt625ZCcgin
2ro/eOkNyUOTBIbuj9CvMnhUYiR61lC1f1IGbrSYYimqBVSjpifVufxtx/I3exRe
ZosTByYp4Xwpb1+WAQIDAQAB
-----END PUBLIC KEY-----
`)
var rsaPublicKey *rsa.PublicKey

func getRsaPublicKey() (*rsa.PublicKey, error) {
	var err error = nil
	if rsaPublicKey == nil {
		rsaPublicKey, err = crypto.ParsePublicKey(publicKey)
	}
	return rsaPublicKey, err
}
func SearchSong(key common.MapType) common.Song {
	searchSong := common.Song{
	}
	keyword := key["keyword"].(string)
	searchSongName := key["name"].(string)
	searchSongName = strings.ToUpper(searchSongName)
	searchArtistsName := key["artistsName"].(string)
	searchArtistsName = strings.ToUpper(searchArtistsName)
	header := make(http.Header, 2)
	header["origin"] = append(header["origin"], "http://music.migu.cn/")
	header["referer"] = append(header["referer"], "http://music.migu.cn/")
	clientRequest := network.ClientRequest{
		Method: http.MethodGet,
		RemoteUrl: "http://m.music.migu.cn/migu/remoting/scr_search_tag?keyword="+keyword+"&type=2&rows=20&pgc=1",
		Host:  "pd.musicapp.migu.cn",
		Proxy: false,
	}
	//fmt.Println(clientRequest.RemoteUrl)
	resp, err := network.Request(&clientRequest)
	if err != nil {
		fmt.Println(err)
		return searchSong
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		return searchSong
	}
	body, err := network.StealResponseBody(resp)
	if err != nil {
		fmt.Println(err)
		return searchSong
	}
	result := utils.ParseJsonV2(body)
	//fmt.Println(utils.ToJson(result))
	var copyrightId = ""
	data, ok := result["musics"].(common.SliceType)
	if ok {
		list:= data
		if ok && len(list) > 0 {
			listLength := len(list)
			for index, matched := range list {
				miguSong, ok := matched.(common.MapType)
				if ok {
					cId, ok := miguSong["copyrightId"].(string)
					if ok {
						singerName, singerNameOk := miguSong["singerName"].(string)
						songName, songNameOk := miguSong["songName"].(string)
						if strings.Contains(songName, "伴奏") && !strings.Contains(searchSongName, "伴奏") {
							continue
						}
						var songNameSores float32 = 0.0
						if songNameOk {
							//songNameKeys := utils.ParseSongNameKeyWord(songName)
							////fmt.Println("songNameKeys:", strings.Join(songNameKeys, "、"))
							//songNameSores = utils.CalMatchScores(searchSongName, songNameKeys)
							songNameSores=utils.CalMatchScoresV2(searchSongName,songName,"songName")
							//fmt.Println("songNameSores:", songNameSores)
						}
						var artistsNameSores float32 = 0.0
						if singerNameOk {
							//artistKeys := utils.ParseSingerKeyWord(singerName)
							////fmt.Println("migu:artistKeys:", strings.Join(artistKeys, "、"))
							//artistsNameSores = utils.CalMatchScores(searchArtistsName, artistKeys)
							artistsNameSores=utils.CalMatchScoresV2(searchArtistsName,singerName,"singerName")
							//fmt.Println("migu:artistsNameSores:", artistsNameSores)
						}
						songMatchScore := songNameSores*0.6 + artistsNameSores*0.4
						//fmt.Println("migu:songMatchScore:", songMatchScore)
						if songMatchScore > searchSong.MatchScore {
							searchSong.MatchScore = songMatchScore
							copyrightId = cId
							searchSong.Name = songName
							searchSong.Artist = singerName
							searchSong.Artist = strings.ReplaceAll(singerName, " ", "")
						}
					}

				}
				if index >= listLength/2 || index > 9 {
					break
				}
			}

		}
	}

	if len(copyrightId) > 0 {
		clientRequest := network.ClientRequest{
			Method:               http.MethodGet,
			RemoteUrl:            "http://music.migu.cn/v3/api/music/audioPlayer/getPlayInfo?dataType=2&" + encrypt("{\"copyrightId\":\""+copyrightId+"\"}"),
			Host:                 "music.migu.cn",
			Header:               header,
			Proxy:                true,
			ForbiddenEncodeQuery: true, //dataType first must
		}
		//fmt.Println(clientRequest.RemoteUrl)
		resp, err := network.Request(&clientRequest)
		if err != nil {
			fmt.Println(err)
			return searchSong
		}
		defer resp.Body.Close()
		body, err = network.StealResponseBody(resp)
		data := utils.ParseJsonV2(body)
		//fmt.Println(data)
		data, ok := data["data"].(common.MapType)
		if ok {
			playInfo, ok := data["sqPlayInfo"].(common.MapType)
			if !ok {
				playInfo, ok = data["hqPlayInfo"].(common.MapType)
				if !ok {
					playInfo, ok = data["bqPlayInfo"].(common.MapType)
				}
			}
			if ok {
				playUrl, ok := playInfo["playUrl"].(string)
				if ok && strings.Index(playUrl, "http") == 0 {
					searchSong.Url = playUrl
					return searchSong
				}
			}
		}
	}
	return searchSong

}

func encrypt(text string) string {
	encryptedData := ""
	//fmt.Println(text)
	text = utils.ToJson(utils.ParseJson(bytes.NewBufferString(text).Bytes()))
	randomBytes, err := utils.GenRandomBytes(32)
	if err != nil {
		fmt.Println(err)
		return encryptedData
	}
	pwd := bytes.NewBufferString(hex.EncodeToString(randomBytes)).Bytes()
	salt, err := utils.GenRandomBytes(8)
	if err != nil {
		fmt.Println(err)
		return encryptedData
	}
	//key = []byte{0xaf, 0xb3, 0xac, 0x50, 0xcd, 0x1d, 0x23, 0x81, 0x58, 0x5f, 0xa7, 0xbc, 0xbd, 0x8c, 0xbe, 0x02, 0x56, 0x0f, 0xad, 0xe7, 0xd1, 0x7e, 0x2e, 0xb1, 0x14, 0x81, 0x6f, 0x27, 0xab, 0x7b, 0x6a, 0x75}
	//iv = []byte{0xfb, 0x10, 0x89, 0xb0, 0x13, 0x32, 0xf2, 0xa7, 0x02, 0x51, 0x49, 0xff, 0xbc, 0x16, 0xf0, 0x40}
	//pwd = bytes.NewBufferString("d8e28215ed6573e0fd5eb8b8ae8062542589e96f669bee6503af003c63cdfbd4").Bytes()
	//salt = []byte{0xde, 0xfc, 0x9f, 0x26, 0x29, 0xdd, 0xec, 0x37}
	key, iv := derive(pwd, salt, 256, 16)
	var data []byte
	data = append(data, bytes.NewBufferString("Salted__").Bytes()...)
	data = append(data, salt...)
	encryptedD := crypto.AesEncryptCBCWithIv(bytes.NewBufferString(text).Bytes(), key, iv)
	data = append(data, encryptedD...)
	dat := base64.StdEncoding.EncodeToString(data)
	var rsaB []byte
	pubKey, err := getRsaPublicKey()
	if err == nil {
		rsaB = crypto.RSAEncryptV2(pwd, pubKey)
	} else {
		rsaB = crypto.RSAEncrypt(pwd, publicKey)
	}
	sec := base64.StdEncoding.EncodeToString(rsaB)
	//fmt.Println("data:", dat)
	//fmt.Println("sec:", sec)
	encryptedData = "data=" + url.QueryEscape(dat)
	encryptedData = encryptedData + "&secKey=" + url.QueryEscape(sec)
	return encryptedData
}
func derive(password []byte, salt []byte, keyLength int, ivSize int) ([]byte, []byte) {
	keySize := keyLength / 8
	repeat := math.Ceil(float64(keySize+ivSize*8) / 32)
	var data []byte
	var lastData []byte
	for i := 0.0; i < repeat; i++ {
		var md5Data []byte
		md5Data = append(md5Data, lastData...)
		md5Data = append(md5Data, password...)
		md5Data = append(md5Data, salt...)
		h := md5.New()
		h.Write(md5Data)
		md5Data = h.Sum(nil)
		data = append(data, md5Data...)
		lastData = md5Data
	}
	return data[:keySize], data[keySize : keySize+ivSize]
}
