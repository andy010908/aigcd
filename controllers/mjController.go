package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "aigcd/core/logger"
	core "aigcd/core/mysql"
	"aigcd/mj"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ImageItem struct {
	ImageID  int64  `json:"image_id"`
	TaskID   string `json:"task_id"`
	ImageUrl string `json:"image_url"`
}

type PromptPreView struct {
	Preview []ImageItem `json:"preview"`
}

type RespPrompt struct {
	Code int           `json:"code"`
	Data PromptPreView `json:"data"`
}

type RespError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type ApplicationCommandOption struct {
	Type        uint8  `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required"`
}

type ApplicationCommand struct {
	ID            string `json:"id,omitempty"`
	ApplicationID string `json:"application_id,omitempty"`
	GuildID       string `json:"guild_id,omitempty"`
	Version       string `json:"version,omitempty"`
	Type          uint8  `json:"type,omitempty"`
	Name          string `json:"name"`
	// NOTE: DefaultPermission will be soon deprecated. Use DefaultMemberPermissions and DMPermission instead.
	DefaultPermission        bool                       `json:"default_permission,omitempty"`
	DefaultMemberPermissions int64                      `json:"default_member_permissions,string,omitempty"`
	DMPermission             bool                       `json:"dm_permission,omitempty"`
	NSFW                     bool                       `json:"nsfw,omitempty"`
	Description              string                     `json:"description,omitempty"`
	Options                  []ApplicationCommandOption `json:"options"`
}

type PromptOption struct {
	Type  int    `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DataPrompt struct {
	Version        string             `json:"version"`
	ID             string             `json:"id"`
	Name           string             `json:"name"`
	Type           int                `json:"type"`
	Options        []PromptOption     `json:"options"`
	ApplicationCmd ApplicationCommand `json:"application_command"`
}

type PromptCommand struct {
	Type          int        `json:"type,omitempty"`
	ApplicationId string     `json:"application_id"`
	GuildId       string     `json:"guild_id"`
	ChannelId     string     `json:"channel_id"`
	SessionId     string     `json:"session_id"`
	Data          DataPrompt `json:"data"`
	Nounce        string     `json:"nounce"`
}

type BodyPrompt struct {
	Prompt string `json:"prompt"`
}

type RespCollectionJSON struct {
	Code      int                 `json:"code"`
	ImageList *[]core.Collections `json:"image_list"`
}

const (
	diffusion_apikey = "A4C5235CB7382F2B"
	api_prompt_url   = "https://discord.com/api/v9/interactions"
	application_id   = "936929561302675456"  //midjourney application id
	guild_id         = "1088383794702196826" //server id
	channel_id       = "1088383794702196829" //channel id
)

//harrison's
//guild_id: 1088383794702196826
//channel_id: 1088383794702196829
//application_id: 1088388251481554954

// var DisCordToken string = "MTA3NzUxOTM0NzIyNTk4NTA0NA.GrSuJH.dC3qH4BlTM25F-tf4-XpwmBoy_cPjeF52osyzI"
var DisCordToken string = "MTA3NzUxOTM0NzIyNTk4NTA0NA.GsHsfY.dyDyrqRhV5SiGL9_aSppQo-jHhKfGLDXbHgGkU"

var session_id string

func CloudUpload(c *gin.Context) {
	mj.UploadFile("aigcd", "harrison0617_acquaint_coffee_shop_in_Paris_with_a_foggy_backgro_eb2da07b-e024-41cb-8bd8-80af728262a3_0.png")
	r := &RespError{}
	r.Code = 0
	r.Msg = "success"
	c.JSON(200, r)
}

func GetCollections(c *gin.Context) {
	apikey := c.Request.Header.Get("apikey")
	if apikey != diffusion_apikey {
		log.Error("GetCollections: apikey error")
		resp := &RespError{}
		resp.Code = 2
		resp.Msg = "apikey error"
		c.JSON(400, resp)
		return
	}
	strCount := c.Request.Header.Get("count")
	if strCount == "" {
		log.Error("GetCollections: count missed")
		resp := &RespError{}
		resp.Code = 2
		resp.Msg = "count missed"
		c.JSON(400, resp)
		return
	}
	count, err := strconv.Atoi(strCount)
	if err != nil {
		log.Error("GetCollections: count error")
		resp := &RespError{}
		resp.Code = 2
		resp.Msg = "count error"
		c.JSON(400, resp)
		return
	}
	strPage := c.Request.Header.Get("page")
	if strPage == "" {
		log.Error("GetCollections: page missed")
		resp := &RespError{}
		resp.Code = 2
		resp.Msg = "count page"
		c.JSON(400, resp)
		return
	}
	page, err := strconv.Atoi(strPage)
	if err != nil {
		log.Error("GetCollections: page error")
		resp := &RespError{}
		resp.Code = 2
		resp.Msg = "page error"
		c.JSON(400, resp)
		return
	}

	db := core.DB.Model(&core.Collections{}).Limit(count).Offset((page - 1) * count)
	var col []core.Collections
	err = db.Order("id DESC").Scan(&col).Error
	if err == nil {
		log.Info("GetCollections", zap.Any("collections", col))
		resp := &RespCollectionJSON{}
		resp.Code = 0
		resp.ImageList = &col
		c.JSON(200, resp)
	} else {
		resp := &RespError{}
		resp.Code = 2
		resp.Msg = "db error"
		c.JSON(400, resp)
	}

}

func SendPrompt(c *gin.Context) {
	apikey := c.Request.Header.Get("apikey")
	if apikey != diffusion_apikey {
		log.Error("SendPrompt: apikey error")
		resp := &RespError{}
		resp.Code = 2
		resp.Msg = "apikey error"
		c.JSON(400, resp)
		return
	}
	promptInput := BodyPrompt{}
	c.BindJSON(&promptInput)
	//log.Info("diffusion", zap.String("prompt", promptInput.Prompt))

	body := PromptCommand{}
	body.Type = 2
	body.ApplicationId = application_id
	body.GuildId = guild_id
	body.ChannelId = channel_id
	body.SessionId = uuidGen()
	body.Nounce = nounceGen()

	//log.Info("diffusion body1")

	//body.Data.Version = "1077969938624553050"
	body.Data.Version = "1237876415471554623"
	body.Data.ID = "938956540159881230"
	body.Data.Name = "imagine"
	body.Data.Type = 1
	body.Data.Options = make([]PromptOption, 1)
	body.Data.Options[0].Type = 3
	body.Data.Options[0].Name = "prompt"
	body.Data.Options[0].Value = promptInput.Prompt

	//log.Info("diffusion body2")

	body.Data.ApplicationCmd.ID = "938956540159881230"
	body.Data.ApplicationCmd.ApplicationID = application_id
	body.Data.ApplicationCmd.Version = "1237876415471554623"
	body.Data.ApplicationCmd.DefaultPermission = true
	//body.Data.ApplicationCmd.DefaultMemberPermissions = ""
	body.Data.ApplicationCmd.Type = 1
	body.Data.ApplicationCmd.NSFW = false
	body.Data.ApplicationCmd.Name = "imagine"
	body.Data.ApplicationCmd.Description = "Create images with Midjourney"
	body.Data.ApplicationCmd.DMPermission = true
	body.Data.ApplicationCmd.Options = make([]ApplicationCommandOption, 1)
	body.Data.ApplicationCmd.Options[0].Type = 3
	body.Data.ApplicationCmd.Options[0].Name = "prompt"
	body.Data.ApplicationCmd.Options[0].Description = "The prompt to imagine"
	body.Data.ApplicationCmd.Options[0].Required = true

	//log.Info("diffusion body3")

	marshalled, err := json.Marshal(body)
	if err != nil {
		log.Error("diffusion", zap.Error(err))
		resp := &RespError{}
		resp.Code = 1
		resp.Msg = "json marshall error"
		c.JSON(500, resp)
		return
	}
	//log.Info("diffusion", zap.String("body marshall", string(marshalled)))

	req, err := http.NewRequest("POST", api_prompt_url, bytes.NewReader(marshalled))
	if err != nil {
		log.Error("diffusion", zap.Error(err))
		resp := &RespError{}
		resp.Code = 1
		resp.Msg = "Create http request error"
		c.JSON(500, resp)
		return
	}
	req.Header.Add("authorization", DisCordToken)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("diffusion", zap.Error(err))
		resp := &RespError{}
		resp.Code = 1
		resp.Msg = "Request mj error"
		c.JSON(500, resp)
		return
	}
	defer resp.Body.Close()
	//log.Info("diffusion http.client")
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("diffusion", zap.Error(err))
		resp := &RespError{}
		resp.Code = 1
		resp.Msg = "io error"
		c.JSON(500, resp)
		return
	}
	log.Info("diffusion", zap.String("imagine resp", string(content)))
	r := &RespError{}
	r.Code = 0
	r.Msg = "success"
	c.JSON(200, r)

}

func nounceGen() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var min int64 = 1000000000000000000
	var max int64 = 9200000000000000000
	num := rand.Int63n(max-min+1) + min
	return strconv.FormatInt(int64(num), 10)

}

func uuidGen() string {
	uuid := uuid.New().String()
	uuidWithoutHyphens := strings.Replace(uuid, "-", "", -1)
	return uuidWithoutHyphens
}
