package gql

import (
	"html/template"
	nethttp "net/http"

	"github.com/go-kratos/kratos/v2/transport/http"
)

var page = template.Must(template.New("voyager").Parse(`<!doctype html>
<html>
  <head>
    <style>
      body {
        height: 100%;
        margin: 0;
        width: 100%;
        overflow: hidden;
      }

      #voyager {
        height: 100vh;
      }
    </style>

    <link
      rel="stylesheet"
      href="https://cdn.jsdelivr.net/npm/graphql-voyager/dist/voyager.css"
    />
    <script src="https://cdn.jsdelivr.net/npm/graphql-voyager/dist/voyager.standalone.js"></script>
  </head>

  <body>
    <div id="voyager">Loading...</div>
    <script type="module">
      const { voyagerIntrospectionQuery: query } = GraphQLVoyager;
      const response = await fetch(
        {{.Endpoint}},
        {
          method: 'post',
          headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ query }),
          credentials: 'omit',
        },
      );
      const introspection = await response.json();

      // Render <Voyager /> into the body.
      GraphQLVoyager.renderVoyager(document.getElementById('voyager'), {
        introspection,
      });
    </script>
  </body>
</html>
`))

func GraphqlVoyager(endpoint string) nethttp.HandlerFunc {
	data := struct {
		Endpoint string
	}{
		Endpoint: endpoint,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=UTF-8")
		if err := page.Execute(w, data); err != nil {
			panic(err)
		}
	}
}
