package common

import (
  "encoding/json"
  "wikiproxy/pkg/common/log"
  "github.com/gomodule/redigo/redis"
)

type Redis struct{
  Host   string
  Logger *log.Logger
}

func NewRedis(host string,logger *log.Logger)*Redis{
   return &Redis{Host: host,Logger: logger}
}

func (r *Redis)GetKey(tokename string)(string,error){
   var token string
   c,err := redis.Dial("tcp",r.Host)
   if err != nil{
     r.Logger.Errorln("dial redis host err: ",err)
     return "",err
   }
   defer c.Close()
   if apptoken,err := c.Do("get",tokename);err != nil{
        r.Logger.Errorln("get key: ",tokename," err: ",err)
        return "",err
     }else{
      if v,ok := apptoken.([]byte);ok{
            json.Unmarshal(v,&token)
      }
   }
   return token,nil
}

func (r *Redis)SetKey(tokename,timeout string,v interface{})error{
   p,err := json.Marshal(v)
   if err != nil{
      r.Logger.Errorln("json marshal err: ",err)
      return err
   }
   c,err := redis.Dial("tcp",r.Host)
   if err != nil{
     r.Logger.Errorln("dial redis host err: ",err)
     return err
   }
   defer c.Close()
   if _,err := c.Do("set",tokename,p,"EX",timeout);err != nil{
      r.Logger.Errorln("set key err: ",err)
      return err
   }
   return nil
}



func (r *Redis)GetKeyTTL(tokename string)(int64,error){
   var ttl int64
   c,err := redis.Dial("tcp",r.Host)
   if err != nil{
     r.Logger.Errorln("dial redis host err: ",err)
     return 0,err
   }
   defer c.Close()
   tokenTTL,err := c.Do("ttl",tokename)
   if err != nil{
      r.Logger.Errorln("get key ttl err: ",err)
      return 0,err
   }else{
    if v,ok := tokenTTL.(int64);ok{
      ttl = v
    }
   }
   return ttl,nil
}


func (r *Redis)DelKey(tokename string)error{
   c,err := redis.Dial("tcp",r.Host)
   if err != nil{
     r.Logger.Errorln("dial redis host err: ",err)
     return err
   }
   defer c.Close()

   _,err = c.Do("del",tokename)
   if err != nil{
      r.Logger.Errorln("del key: ",tokename," err: ",err)
      return err
   } 
   return nil
}



func (r *Redis)CacheRe(name string)interface{}{
   c,err := redis.Dial("tcp",r.Host)
   if err != nil{
     r.Logger.Errorln("dial redis host err: ",err)
     return nil
   }
   defer c.Close()

   data,err := c.Do("get",name)
   if err != nil{
      r.Logger.Errorln(err)
      return nil
   }
   return data
}

func (r *Redis)CacheSet(name string,v interface{})error{
   p,err := json.Marshal(v)
   if err != nil{
     r.Logger.Errorln(err)
   }
   c,err := redis.Dial("tcp",r.Host)
   if err != nil{
     r.Logger.Errorln("dial redis host err: ",err)
     return err
   }
   defer c.Close()
   if _,err := c.Do("set",name,p);err != nil{
      return err
   }
   r.Logger.Debugln("set ",name,"to cache server successed.")
   return nil
}






