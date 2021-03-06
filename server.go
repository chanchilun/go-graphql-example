package main

import (
	"github.com/OscarYuen/go-graphql-example/conf"
	"github.com/OscarYuen/go-graphql-example/model"
	"github.com/OscarYuen/go-graphql-example/schema"
	"log"
	"net/http"

	"github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
)

var gschema *graphql.Schema

func init() {
	gschema = graphql.MustParseSchema(schema.GetRootSchema(), &schema.Resolver{})
}

func main() {
	db := conf.ConnectDB()
	schema.SetDatabase(db)
	// Migrate the schema
	db.AutoMigrate(&model.User{})

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))

	http.Handle("/query", &relay.Handler{Schema: gschema})

	log.Fatal(http.ListenAndServe(":3000", nil))
}

var page = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.css" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/1.1.0/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react-dom.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.js"></script>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function graphQLFetcher(graphQLParams) {
				return fetch("/query", {
					method: "post",
					body: JSON.stringify(graphQLParams),
					credentials: "include",
				}).then(function (response) {
					return response.text();
				}).then(function (responseBody) {
					try {
						return JSON.parse(responseBody);
					} catch (error) {
						return responseBody;
					}
				});
			}
			ReactDOM.render(
				React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
				document.getElementById("graphiql")
			);
		</script>
	</body>
</html>
`)
