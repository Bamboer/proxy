package common

import (
  "io"
  "fmt"
  "time"
  "bytes"
  "strconv"
  "net/url"
  "net/http"
  "io/ioutil"
  "mime/multipart"
  "encoding/json"
  "wxcrm/pkg/common"
  "wxcrm/pkg/common/log"

)



//企业微信access token 返回data struct
type WXRT struct{
  Errcode      int          `json:"errcode,required"`
  Errmsg       string       `json:"errmsq,required"`
  Access_token string       `json:"access_token,required"`
  Expires_in   int          `json:"expires_in,required"`
}

//企业微信当前使用用户ID struct
type AccessUser struct{
  Errcode int          `json:"errcode,required"`
  Errmsg  string       `json:"errmsq,required"`
  UserId  string       `json:"userid,required"`
  DeviceId string      `json:"deviceid,required"`
}


type WX struct{
  Redis   *common.Redis
  Logger  *log.Logger
  CorpID  string 
  AppSecret string
  ContactSecret string  
  AgentId       string 
  RecycleEventsToken    string  
  RecycleEventsAesKey  string

  ContactToken string //WX 通讯录token
  Token        string
  JSTicket     string 
  TokenSetTime int64
  TicketSetTime int64
}

const (
  WXTokenName   = "wxtoken"
  JSTicketName  = "jsticket"
  WXContactTokenName = "contactoken"
  // YYTokenName   = "yytoken"
  TokenTTL    = "7200"
)

func NewWX(cfg *common.Opts,redis *common.Redis,logger *log.Logger)*WX{
  return &WX{
    Redis: redis,
    Logger: logger,
    CorpID: cfg.WXCorpId,
    AppSecret: cfg.WXAppSecret,
    ContactSecret: cfg.WXContactSecret,
    AgentId:  cfg.WXAgentId,
    RecycleEventsToken: cfg.WXRecycleEventsToken,   
    RecycleEventsAesKey: cfg.WXRecycleEventsAesKey,  
  }
}



// 获取应用access token
func (wx *WX)GetWXToken()(string,error){
  tnow := time.Now().Unix()
  if tnow - wx.TokenSetTime < -10 && wx.Token != ""{
      // wx.Logger.Debugln("WXToken set time is not timeout and return this token value.")
      // wx.Logger.Debugln("WX token: ",wx.Token)
      return wx.Token,nil
  }    

  t1,err := wx.Redis.GetKeyTTL(WXTokenName)
  // wx.Logger.Debugln("Get token :",WXTokenName,"ttl: ",t1)
  if err != nil || t1 == -2{
    wxrt := &WXRT{}
    client := &http.Client{}

    // https://qyapi.weixin.qq.com/cgi-bin/gettoken

    uri,_ := url.Parse("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + wx.CorpID +"&corpsecret=" +wx.AppSecret )
    req,err := http.NewRequest("GET", uri.String(), nil)
    if err != nil{
      wx.Logger.Errorln(err)
      return "",err
    }

    wx.Logger.Debugln("GET WX Token url: ",req.URL.String())
  
    resq,err := client.Do(req)
    defer resq.Body.Close()
  
    b,err := ioutil.ReadAll(resq.Body)
    if err != nil {
      wx.Logger.Errorln(err)
      return "",err
    }
    if err := json.Unmarshal(b,wxrt);err != nil{
      wx.Logger.Errorln(err)
      return "",err
    }
    if wxrt.Errcode == 0 {
      if err := wx.Redis.SetKey(WXTokenName,strconv.Itoa(wxrt.Expires_in),wxrt.Access_token);err != nil{
         wx.Logger.Errorln(err)
      }

      tnow := time.Now().Unix()
      wx.Token = wxrt.Access_token
      wx.TokenSetTime = tnow + int64(wxrt.Expires_in)

      return wxrt.Access_token,nil
    }else{
      wx.Logger.Errorln(wxrt.Errmsg)
      return "",fmt.Errorf(wxrt.Errmsg)
   }
  }else{
    token,err := wx.Redis.GetKey(WXTokenName)
    if err != nil{
      wx.Logger.Errorln(err)
       return "",err
    }
    // wx.Logger.Debugln("WX token: ",token)
    return token,nil
  }
}

//获取企业微信当前用户
func (wx *WX)GetWXAccessUser(code string)(string,error){
  //https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=ACCESS_TOKEN&code=CODE
  accessuser := &AccessUser{}
  client := &http.Client{}

  token,err := wx.GetWXToken()
  if err != nil{
    wx.Logger.Errorln(err)
    return "",err
  }
  uri,_ := url.Parse("https://qyapi.weixin.qq.com"+"/cgi-bin/user/getuserinfo?access_token="+token+"&code="+code+"&debug=1")
  req,err := http.NewRequest("GET", uri.String(), nil)
  if err != nil{
     wx.Logger.Errorln(err)
     return "",err
  }

  resq,err := client.Do(req)
  if err != nil{
     wx.Logger.Errorln(err)
     return "",err
  }

  defer resq.Body.Close()

  b,err := ioutil.ReadAll(resq.Body)
  if err != nil {
     wx.Logger.Errorln(err)
     return "",err
  }
  if err := json.Unmarshal(b,accessuser);err != nil{
     wx.Logger.Errorln(err)
     return "",err
  }
  if accessuser.Errcode == 0 {
     wx.Logger.Debugln("Current userID: ",accessuser.UserId)
     return accessuser.UserId,nil
  }else{
     wx.Logger.Errorln("Get userid Error code: ",accessuser.Errcode," msg: ",accessuser)
     return "",fmt.Errorf("get wx exuserid error.")
  }
}

//读取成员详细信息

func (wx *WX)GetWXAccessUserDetail(userid string)(*WXUserDetail,error){
  //GET https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=ACCESS_TOKEN&userid=USERID
  userdetail := &WXUserDetail{}
  client := &http.Client{}

  token,err := wx.GetWXToken()
  if err != nil{
    wx.Logger.Errorln(err)
    return nil,err
  }
  uri,_ := url.Parse("https://qyapi.weixin.qq.com"+"/cgi-bin/user/get?access_token="+token+"&userid="+userid)
  req,err := http.NewRequest("GET", uri.String(), nil)
  if err != nil{
     wx.Logger.Errorln(err)
     return nil,err
  }

  resq,err := client.Do(req)
  if err != nil{
     wx.Logger.Errorln(err)
     return nil,err
  }

  defer resq.Body.Close()

  b,err := ioutil.ReadAll(resq.Body)
  if err != nil {
     wx.Logger.Errorln(err)
     return nil,err
  }
  if err := json.Unmarshal(b,userdetail);err != nil{
     wx.Logger.Errorln(err)
     return nil,err
  }
  if userdetail.Errcode == 0 {
     wx.Logger.Debugln("Current user name: ",userdetail.Name)
     return userdetail,nil
  }else{
     wx.Logger.Errorln(userdetail.Errmsg)
     return nil,fmt.Errorf(userdetail.Errmsg)
  }
}

