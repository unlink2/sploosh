package restapi

import (
    "github.com/emicklei/go-restful"
    "config"
    "errors"
)

type IrcBridgeResponse struct {
  Message string
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

  // get whitelisted arrays
  clientids := config.Globalcfg.Section("clientids")
  clientidKeys := clientids.Keys()

  var canExecute = false

  for _, key := range clientidKeys {
    if key.String() == clientid {
      canExecute = true
      break;
    }
  }

  if !canExecute && len(clientidKeys) > 0 {
    response.WriteError(418, errors.New("Invalid client-id"))
    return
  }


  response.WriteEntity(IrcBridgeResponse{})
}
