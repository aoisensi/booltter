package main

import (
    "os"
    "fmt"
    "flag"
    "github.com/golang/glog"
    "github.com/ChimeraCoder/anaconda"
    "github.com/gin-gonic/gin"
    "github.com/garyburd/go-oauth/oauth"
    "github.com/gin-gonic/contrib/sessions"
)

func main() {
    flag.Parse()
    r := gin.Default()
    store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
    r.Use(sessions.Sessions("session", store))
    
    r.GET("/", func(ctx *gin.Context) {
        session := sessions.Default(ctx)
        id := session.Get("id")
        if id == nil {
            ctx.String(200, "not yet login")
            return
        }
        ctx.String(200, "ur twitter id is %v", id.(int64))
        
    })
    
    
    r.GET("/signin", func(ctx *gin.Context) {
        session := sessions.Default(ctx)
        url, cre, err := anaconda.AuthorizationURL("")
        
        if err != nil {
            glog.Error(err)
            ctx.String(500, "signin failed.")
            return
        }
        
        session.Set("credentials_token", cre.Token)
        session.Set("credentials_secret", cre.Secret)
        if err := session.Save(); err != nil {
            glog.Warning(err)
        }
        
        ctx.Redirect(303, url)
        return
    })
    
    r.GET("/signin/callback", func(ctx *gin.Context) {
        session := sessions.Default(ctx)
        token, ok := session.Get("credentials_token").(string)
        if !ok {
            ctx.String(500, "callback failed.")
            return
        }
        secret, ok := session.Get("credentials_secret").(string)
        if !ok {
            ctx.String(500, "callback failed.")
            return
        }
        
        cre := &oauth.Credentials{
            Token: token,
            Secret: secret,
        }
        verifier := ctx.Request.URL.Query().Get("oauth_verifier")
        cred, _, err := anaconda.GetCredentials(cre, verifier)
        defer func() {
            session.Delete("credentials")
            if err := session.Save(); err != nil {
                glog.Warning(err)
            }
        }()
        if err != nil {
            ctx.String(500, "callback failed.")
            return
        }
        
        user, err := findOrCreateUserFromToken(cred.Token, cred.Secret)
        if err != nil {
            ctx.String(500, "callback failed.")
            return
        }
        
        session.Set("id", user.ID)
        ctx.Redirect(303, "/")
        
        
    })
    bind := fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
    r.Run(bind)
}

func signin(session sessions.Session, token *oauth.Credentials) error {
    user, err := findOrCreateUserFromToken(token.Token, token.Secret)
    if err != nil {
        return err
    }
    session.Set("id", user.ID)
    session.Save()
    return nil
}