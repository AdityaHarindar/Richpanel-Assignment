# Richpanel-Assignment

**Setup Instructions**
1. Install Go 1.24.4
2. `go mod download` from project root
3. `go run main.go` from project root to start the API server

**Design:**
Employs a multi layered approach - project split into model/store/transport packages

**_model_**: contains the data model for blog posts

    type Post struct {
        ID        int       `json:"id"`
        Title     string    `json:"title"`
        Author    string    `json:"author"`
        Body      string    `json:"body"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
    }

**_store_**: contains all code related to different data stores, hidden behind the store/cache interface, allows for switching to different data stores without any addition modification to other pacakges as long as they satsify the interface spec

**_transport_**: contains all related code for transport being used, HTTP in this case, can be switched out to grpc, etc 


Additional improvement:
1. Unit/Integration tests with mocks, can be run using `go test`
2. Benchmarks for handler level methods, using `go bench`
3. Documentation can be generated using `go doc` [server](https://go.dev/blog/godoc) based on complying formatted comments 