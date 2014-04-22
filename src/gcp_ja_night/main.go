package gcp_ja_night

import (
	"net/http"
	"github.com/crhym3/go-endpoints/endpoints"
	"time"
	"github.com/mjibson/goon"
	"appengine/datastore"
	"fmt"
)

type Greeting struct {
	Key     int64         `json:"id" datastore:"-" goon:"id"`
	Author  string         `json:"author"`
	Content string         `json:"content" datastore:",noindex"`
	Date    time.Time      `json:"date"`
}
type GreetingList struct {
	Items []*Greeting       `json:"items"`
}

type GreetingListReq struct {
	Limit int         `json:"limit" endpoints:"d=10,min=1,max=100,desc=The limit search result count"`
}

type GreetingGetReq struct {
	Key int         `json:"key" endpoints:"req,required"`
}


type GreetingService struct {}

func (gs *GreetingService) Get(r *http.Request, req *GreetingGetReq, resp *Greeting) error  {
	//c := endpoints.NewContext(r)

	g := goon.NewGoon(r)
	resp.Key = int64(req.Key)

	if err := g.Get(resp);err != nil {
		if err == datastore.ErrNoSuchEntity {
			return endpoints.NewNotFoundError("Not Found")
		}
		return endpoints.NewInternalServerError(err.Error())
	}

	return nil
}

func (gs *GreetingService) List(r *http.Request, _ *GreetingListReq, resp *GreetingList) error {
	//endpoints.Contextを使うと、Google OAuthからのユーザが取得できます。
	//c := endpoints.NewContext(r)

	g := goon.NewGoon(r)

	var items []*Greeting

	if _, err := g.GetAll(datastore.NewQuery(g.Kind(&Greeting{})), &items); err != nil {
		return endpoints.NewInternalServerError(err.Error())
	}

	resp.Items = items

	return nil
}

func (gs *GreetingService) Insert(r *http.Request, req *Greeting, resp *Greeting) error {
	c := endpoints.NewContext(r)
	g := goon.NewGoon(r)

	c.Debugf("%s" , req.Key)

	if req.Key != 0 {
		if err := g.Get(&Greeting{Key : req.Key}); err == nil {
			return endpoints.NewConflictError(fmt.Sprintln("Found %d", req.Key))
		}
	}

	if _, err := g.Put(req); err != nil {
		return endpoints.NewInternalServerError(err.Error())
	}

	resp.Key = req.Key
	resp.Author = req.Author
	resp.Content = req.Content
	resp.Date = req.Date

	return nil

}

func init() {
	//サービスのインスタンスを作って RegisterServiceを呼び出し
	greetService := &GreetingService{}
	api, err := endpoints.RegisterService(greetService, "greeting", "v1", "Greetings API", true)
	if err != nil {
		panic(err.Error())
	}

	info := api.MethodByName("Get").Info()
	info.Name, info.HttpMethod, info.Path, info.Desc = "greets.get", "GET", "{key}", "Get a Greet."

	info = api.MethodByName("List").Info()
	info.Name, info.HttpMethod, info.Path, info.Desc = "greets.list", "GET", "list", "Get Greets."

	info = api.MethodByName("Insert").Info()
	info.Name, info.HttpMethod, info.Path, info.Desc = "greets.insert", "POST", "insert", "Insert new Greet."


	endpoints.HandleHttp()
}
