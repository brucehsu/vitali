package resources

import (
    "fmt"
    "strconv"
    "net/http"
    "github.com/lunastorm/vitali"
    "github.com/lunastorm/vitali/example/util"
)

type Progress struct {
    vitali.Ctx
    vitali.Perm `*:"AUTHED"`
    ChanMap *util.ChanMap
}

func (c *Progress) Get() interface{} {
    clientClosec := c.ResponseWriter.(http.CloseNotifier).CloseNotify()
    progressc := c.ChanMap.Get(c.Username, c.PathParam("slide"), c.Request.RemoteAddr)
    defer c.ChanMap.Remove(c.Username, c.PathParam("slide"), c.Request.RemoteAddr)
    for ;; {
        select {
        case progress := <-progressc:
            return fmt.Sprintf("%d", progress)
        case <- clientClosec:
            return c.ClientGone()
        }
    }
}

func (c *Progress) Post() interface{} {
    page, err := strconv.ParseUint(c.Param("page"), 10, 32)
    if err != nil {
        return c.BadRequest("bad page")
    }
    c.ChanMap.Broadcast(c.Username, c.PathParam("slide"), int(page))
    return c.NoContent()
}
