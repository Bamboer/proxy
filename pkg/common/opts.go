package common
import (
        "gopkg.in/ini.v1"
        "os"
        "wxcrm/pkg/common/log"
        "github.com/jpillora/opts"
)

var (
  version = "1.0.0"
)

type Opts struct{
        Port               int      `opts:"help=listening port, default=80, env=PORT"`
        ServerName         string   `opts:"help=Access server dns name"`
        RedisHost          string   `opts:"help=redis server addresss, default=x.x.x.x:6379"`
        MysqlHost          string   `opts:"help=mysql server addresss, default=x.x.x.x:3306"`
        MysqlUser          string   `opts:"help=mysql server username, default=proxy"`
        MysqlPassword      string   `opts:"help=mysql server user password, default=proxy123"`
        DatabaseName       string   `opts:"help=mysql server database name,default=proxy"`
        WXCorpId           string   `opts:"help=wx corp id, default=xxx"`
        WXAppSecret        string   `opts:"help=wx app secret, default=xxx"`
        WXAgentId          string   `opts:"help=wx app agent id, default=xxx"`
        LogFile            string   `opts:"help=application log file, default=proxy.log"`
        LogLevel           string   `opts:"help=application log level, default=error"`
}

//Configuration fields
type Obj struct {
  Mysql  Mysql
  Redis  Redisc
  WX     WX
  Port   int
  ServerName    string
  LogFile string
  LogLevel string

}

type Mysql struct{
  User   string
  Password string
  Host    string
  DBname  string
}

type Redisc struct{
  Host   string
}

type WX struct{
  CorpId  string
  AppSecret string
  AgentId      string 
}



func NewLogger(logfile,LogLevel string)*log.Logger{
    var Logger *log.Logger
    file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
            panic("Failed to open configuration file.")
    }
    switch LogLevel{
    case "debug":
            Logger = log.New(file,"",log.Ldate|log.Ltime|log.Lshortfile,0)
    case "info":
            Logger = log.New(file,"",log.Ldate|log.Ltime|log.Lshortfile,1)
    case "warning":
            Logger = log.New(file,"",log.Ldate|log.Ltime|log.Lshortfile,2)
    case "error":
            Logger = log.New(file,"",log.Ldate|log.Ltime|log.Lshortfile,3)
    case "fatal":
            Logger = log.New(file,"",log.Ldate|log.Ltime|log.Lshortfile,4)
    default:
            Logger = log.New(file,"",log.Ldate|log.Ltime|log.Lshortfile,1)
    }
    return Logger
}



func NewOpts()*Opts{
  opt := Opts{}
  opts.New(&opt).Name("proxy").PkgRepo().Version(version).Parse()
  ValidateOpts(&opt)
  return &opt
}


func ValidateOpts(c *Opts){
        var cfg = &Obj{}
        if c.Port == 0 {
        	c.Port = 80
        }
        if c.ServerName == "" {
        	c.ServerName = ""
        }
         
        if c.RedisHost == ""{
                c.RedisHost = "10.12.14.36:6379"
        }
        if c.LogFile == ""{
                c.LogFile = "proxy.log"
        }
        if c.LogLevel == ""{
                c.LogLevel = "error"
        }
        if c.MysqlHost == ""{
                c.MysqlHost = "10.12.14.36:3306"
        }
        if c.MysqlUser == ""{
                c.MysqlUser = "proxy"
        }
        if c.MysqlPassword == ""{
                c.MysqlPassword = "proxy123"
        }
        if c.DatabaseName == ""{
                c.DatabaseName = "proxy"
        }
        if c.WXCorpId == ""{
                c.WXCorpId = "xxx"
        }
        if c.WXAppSecret == ""{
                c.WXAppSecret =  "xxx"
        }    

        if c.WXAgentId == ""{
                c.WXAgentId = "100001xx"
        }  
}
