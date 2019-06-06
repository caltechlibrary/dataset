
# Developer notes

## Requirements

Compiler requirements

+ go1.12 or better

Package assumptions

+ blevesearch v0.7.0 or better

## Recommend

+ Google Sheets support docs see https://developers.google.com/sheets/api/quickstart/go
+ Setup Credentials Wizard: https://console.developers.google.com/start/api?id=sheets.googleapis.com

Run the following to confirm setup

```shell
    go get -u google.golang.org/api/sheets/v4
    go get -u golang.org/x/oauth2/...
```

Per the docs you can test a working connect with the _quickstart.go_ program below.

```go
    package main
    
    import (
      "encoding/json"
      "fmt"
      "io/ioutil"
      "log"
      "net/http"
      "net/url"
      "os"
      "os/user"
      "path/filepath"
    
      "golang.org/x/net/context"
      "golang.org/x/oauth2"
      "golang.org/x/oauth2/google"
      "google.golang.org/api/sheets/v4"
    )
    
    // getClient uses a Context and Config to retrieve a Token
    // then generate a Client. It returns the generated Client.
    func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
      cacheFile, err := tokenCacheFile()
      if err != nil {
        log.Fatalf("Unable to get path to cached credential file. %v", err)
      }
      tok, err := tokenFromFile(cacheFile)
      if err != nil {
        tok = getTokenFromWeb(config)
        saveToken(cacheFile, tok)
      }
      return config.Client(ctx, tok)
    }
    
    // getTokenFromWeb uses Config to request a Token.
    // It returns the retrieved Token.
    func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
      authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
      fmt.Printf("Go to the following link in your browser then type the "+
        "authorization code: \n%v\n", authURL)
    
      var code string
      if _, err := fmt.Scan(&code); err != nil {
        log.Fatalf("Unable to read authorization code %v", err)
      }
    
      tok, err := config.Exchange(oauth2.NoContext, code)
      if err != nil {
        log.Fatalf("Unable to retrieve token from web %v", err)
      }
      return tok
    }
    
    // tokenCacheFile generates credential file path/filename.
    // It returns the generated credential path/filename.
    func tokenCacheFile() (string, error) {
      usr, err := user.Current()
      if err != nil {
        return "", err
      }
      tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
      os.MkdirAll(tokenCacheDir, 0700)
      return filepath.Join(tokenCacheDir,
        url.QueryEscape("sheets.googleapis.com-go-quickstart.json")), err
    }
    
    // tokenFromFile retrieves a Token from a given file path.
    // It returns the retrieved Token and any read error encountered.
    func tokenFromFile(file string) (*oauth2.Token, error) {
      f, err := os.Open(file)
      if err != nil {
        return nil, err
      }
      t := &oauth2.Token{}
      err = json.NewDecoder(f).Decode(t)
      defer f.Close()
      return t, err
    }
    
    // saveToken uses a file path to create a file and store the
    // token in it.
    func saveToken(file string, token *oauth2.Token) {
      fmt.Printf("Saving credential file to: %s\n", file)
      f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
      if err != nil {
        log.Fatalf("Unable to cache oauth token: %v", err)
      }
      defer f.Close()
      json.NewEncoder(f).Encode(token)
    }
    
    func main() {
      ctx := context.Background()
    
      b, err := ioutil.ReadFile("client_secret.json")
      if err != nil {
        log.Fatalf("Unable to read client secret file: %v", err)
      }
    
      // If modifying these scopes, delete your previously saved credentials
      // at ~/.credentials/sheets.googleapis.com-go-quickstart.json
      config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
      if err != nil {
        log.Fatalf("Unable to parse client secret file to config: %v", err)
      }
      client := getClient(ctx, config)
    
      srv, err := sheets.New(client)
      if err != nil {
        log.Fatalf("Unable to retrieve Sheets Client %v", err)
      }
    
      // Prints the names and majors of students in a sample spreadsheet:
      // https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
      spreadsheetId := "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"
      readRange := "Class Data!A2:E"
      resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
      if err != nil {
        log.Fatalf("Unable to retrieve data from sheet. %v", err)
      }
    
      if len(resp.Values) > 0 {
        fmt.Println("Name, Major:")
        for _, row := range resp.Values {
          // Print columns A and E, which correspond to indices 0 and 4.
          fmt.Printf("%s, %s\n", row[0], row[4])
        }
      } else {
        fmt.Print("No data found.")
      }
    }
```

## Using the _dataset_ package

+ create/initialize collection
+ create a JSON document in a collection
+ read a JSON document
+ update a JSON document
+ delete a JSON document

```go
    package main
    
    import (
        "github.com/caltechlibrary/dataset"
        "log"
    )
    
    func main() {
        // Create a collection "mystuff" inside the directory called demo
        collection, err := dataset.InitCollection("demo/mystuff.ds", 
                           dataset.PAIRTREE_LAYOUT)
        if err != nil {
            log.Fatalf("%s", err)
        }
        defer collection.Close()
        // Create a JSON document
        docName := "freda.json"
        document := map[string]interface{}{
            "name":  "freda",
            "email": "freda@inverness.example.org",
        }
        if err := collection.Create(docName, document); err != nil {
            log.Fatalf("%s", err)
        }
        // Read a JSON document
        if err := collection.Read(docName, document); err != nil {
            log.Fatalf("%s", err)
        }
        // Update a JSON document
        document["email"] = "freda@zbs.example.org"
        if err := collection.Update(docName, document); err != nil {
            log.Fatalf("%s", err)
        }
        // Delete a JSON document
        if err := collection.Delete(docName); err != nil {
            log.Fatalf("%s", err)
        }
    }
```


## package requirements

_dataset_ is built on both Golang's standard packages, Caltech Library 
packages and a few 3rd party packages.  At this has not been necessary 
to vendor any packages assuming you're building from the master branch.

## Caltech Library packages

+ [github.com/caltechlibrary/dotpath](https://github.com/caltechlibrary/dotpath)
    + provides dot path style notation to reach into JSON objects
+ [github.com/caltechlibrary/storage](github.com/caltechlibrary/storage)
    + provides a unified storage interaction supporting local disc and AWS S3 storage
+ [github.com/caltechlibrary/tmplfn](https://github.com/caltechlibrary/tmplfn)
    + provides additional template functionality used to format web search results
    + provides a filter engine leveraging the pipeline notation in Go's text templates


## 3rd party packages

+ [Markdown packages] - used to support rendering Markdown embedded in JSON objects
    + [github.com/microcosm-cc/bluemonday](https://github.com/microcosm-cc/bluemonday)
    + [github.com/russross/blackfriday](https://github.com/russross/blackfriday)
+ Migrating to [go-cloud](https://github.com/google/go-cloud) from aws-sdk and Google's Go SDK

