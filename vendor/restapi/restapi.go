package restapi

import (
    "github.com/emicklei/go-restful"
    "config"
    "errors"
    "bot"
    "discord"
    "strings"
)

type IrcBridgeResponse struct {
  Res bot.ResponseWrapper
  Success bool
}

func New() {
  restful.Add(newIRCBridgeService())
}

func newIRCBridgeService() *restful.WebService {
  service := new(restful.WebService)

  service.
    Path("/ircbridge").
    Consumes(restful.MIME_JSON).
    Produces(restful.MIME_JSON)

  service.Route(service.GET("/exec").To(Execute))

  return service
}

func Execute(request *restful.Request, response *restful.Response) {
  response.AddHeader("Access-Control-Allow-Origin", "*")
  response.AddHeader("Access-Control-Allow-Credentials", "true")
  response.AddHeader("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
  response.AddHeader("Access-Control-Allow-Headers", "x-extension-jwt, Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

  clientid := request.QueryParameter("clientid")
  content := request.QueryParameter("content")
  guildid := request.QueryParameter("guildid")
  channelid := request.QueryParameter("channelid")
  authorid := request.QueryParameter("authorid")

  clientsecret := request.QueryParameter("clientsecret")

  if authorid == "" {
    authorid = clientid
  }

  // get whitelisted arrays
  clientids := config.Globalcfg.Section("clientids")
  clientsecrets := config.Globalcfg.Section("clientsecrets")
  clientidKeys := clientids.Keys()

  var success bool
  var res bot.ResponseWrapper

  var canExecute = false

  for _, key := range clientidKeys {
    if key.String() == clientid {
      if clientsecrets.HasKey(key.Name()) {
        if clientsecrets.Key(key.Name()).String() == clientsecret {
          canExecute = true
          break;
        }
      }
    }
  }

  if !canExecute && len(clientidKeys) > 0 {
    response.WriteError(418, errors.New("Invalid client-id or client-secret"))
    return
  }

  // parse stuff here
  for _, command := range bot.Commands {
    for _, name := range command.GetNames() {
      if strings.HasPrefix(content, name) {

        success, res = command.Execute(bot.MessageWrapper{S: nil, M: nil,
          Content: content,

          DChannelID: channelid,
          DGuildID: guildid,
          DAuthorID: authorid,
        })
      }
    }
  }

  if res.Sound != "" && channelid != "" && authorid != "" && guildid != "" {
    go discord.PlayDiscordSound(res, channelid, authorid, nil)
  }

  response.WriteEntity(IrcBridgeResponse{Success: success, Res: res})
}
