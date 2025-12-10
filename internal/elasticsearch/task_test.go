package elasticsearch_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/testcontainers/testcontainers-go"
	elasticsearchtest "github.com/testcontainers/testcontainers-go/modules/elasticsearch"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	elasticsearchtask "github.com/MarioCarrion/todo-api-microservice-example/internal/elasticsearch"
)

const (
	// dockerImage must match the docker image listed in `compose.yml`.
	dockerImage = "elasticsearch:7.17.9"

	// defaultIndexTimeount is timeout to wait for indexing operations.
	defaultIndexTimeount = 5 * time.Second
)

func TestMain(m *testing.M) {
	client := setupClient()
	if client.err != nil {
		panic(fmt.Sprintf("Failed to set up Elasticsearch client: %v", client.err))
	}

	code := m.Run()

	if err := client.Teardown(); err != nil {
		panic(fmt.Sprintf("Failed to close Elasticsearch client: %v", err))
	}

	os.Exit(code)
}

func TestTask_All(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	client := setupClient()
	if client.err != nil {
		t.Fatalf("Failed to setupClient: %v", client.err)
	}

	// Create task repository
	taskRepo := elasticsearchtask.NewTask(client.elasticsearch)

	// Test Index method
	now := time.Now()
	task := internal.Task{
		ID:          "test-123",
		Description: "Test task for elasticsearch",
		Priority:    internal.PriorityHigh.Pointer(),
		IsDone:      false,
		Dates: &internal.Dates{
			Start: &now,
			Due:   &now,
		},
	}

	//- Testing `Index` method
	if err := taskRepo.Index(ctx, task); err != nil {
		t.Fatalf("Failed to index task: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	//- Testing `Search` method
	results, err := taskRepo.Search(ctx, internal.SearchParams{
		Description: internal.ValueToPointer("Test"),
		From:        0,
		Size:        10,
	})
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}

	if len(results.Tasks) == 0 {
		t.Error("Expected to find indexed task")
	}

	//- Testing `Delete` method
	if err = taskRepo.Delete(ctx, task.ID); err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}

	time.Sleep(1 * time.Second)

	results, err = taskRepo.Search(ctx, internal.SearchParams{
		Description: internal.ValueToPointer("Test"),
		From:        0,
		Size:        10,
	})
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}

	if len(results.Tasks) != 0 {
		t.Fatalf("Expected NOT to find indexed task")
	}
}

//-

var setupClient = sync.OnceValue(func() ElasticsearchClient { //nolint: gochecknoglobals
	var res ElasticsearchClient

	ctx := context.Background()

	container, err := elasticsearchtest.Run(ctx, dockerImage)
	if err != nil {
		res.err = fmt.Errorf("failed to start elasticsearch container: %w", err)

		return res
	}

	res.container = container

	//- Create index

	hclient := http.Client{
		Timeout: defaultIndexTimeount,
	}

	u, err := url.Parse(container.Settings.Address)
	if err != nil {
		res.err = fmt.Errorf("failed to parse container address: %w", err)

		return res
	}

	u.Path = "/tasks"

	// Make sure the this payload is the same as the one used in `compose.yml`.
	body := `{"mappings":{"properties":{"id":{"type":"keyword"},"description":{"type":"text"}}}}`

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), strings.NewReader(body))
	if err != nil {
		res.err = fmt.Errorf("failed to parse container address: %w", err)

		return res
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := hclient.Do(req)
	if err != nil {
		res.err = fmt.Errorf("failed to do http request: %w", err)

		return res
	}

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		res.err = fmt.Errorf("failed to do read all response: %w", err)

		return res
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		res.err = fmt.Errorf("failed to create index, status is not 200: %d - %s", resp.StatusCode, string(bodyResp))

		return res
	}

	//- end of creating index

	cfg := elasticsearch.Config{
		Addresses: []string{container.Settings.Address},
	}
	client, err := elasticsearch.NewClient(cfg)

	if err != nil {
		res.err = fmt.Errorf("failed to create elasticsearch client: %w", err)

		return res
	}

	res.elasticsearch = client

	return res
})

type ElasticsearchClient struct {
	elasticsearch *elasticsearch.Client
	container     *elasticsearchtest.ElasticsearchContainer
	err           error
}

func (r *ElasticsearchClient) Teardown() error {
	if r.container != nil {
		if err := testcontainers.TerminateContainer(r.container); err != nil {
			return fmt.Errorf("failed to terminate container: %w", err)
		}
	}

	if r.err != nil {
		return r.err
	}

	return nil
}
