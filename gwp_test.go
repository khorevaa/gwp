package gwp

import (
	"context"
	"errors"
	"fmt"
	"github.com/dalmarcogd/gwp/pkg/worker"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	s := New()
	if s.Configs()["port"] != defaultConfig["port"] {
		t.Errorf("Port is different of default port %d != %d", s.Configs()["port"], defaultConfig["port"])
	}
	if s.Configs()["host"] != defaultConfig["host"] {
		t.Errorf("Host is different of default host %s != %s", s.Configs()["host"], defaultConfig["host"])
	}
	if s.Configs()["basePath"] != defaultConfig["basePath"] {
		t.Errorf("BasePath is different of default basePath %s != %s", s.Configs()["basePath"], defaultConfig["basePath"])
	}
	if s.Configs()["stats"] != defaultConfig["stats"] {
		t.Errorf("Stats is different of default stats %t != %t", s.Configs()["stats"], defaultConfig["stats"])
	}
	if s.Configs()["healthCheck"] != defaultConfig["healthCheck"] {
		t.Errorf("HealthCheck is different of default healthCheck %t != %t", s.Configs()["healthCheck"], defaultConfig["healthCheck"])
	}
	if s.Configs()["debugPprof"] != defaultConfig["debugPprof"] {
		t.Errorf("DebugPprof is different of default debugPprof %t != %t", s.Configs()["debugPprof"], defaultConfig["debugPprof"])
	}
}

func TestNewWithConfig(t *testing.T) {
	config := map[string]interface{}{
		"port":        8002,
		"host":        "google",
		"basePath":    "/workers/test",
		"stats":       true,
		"healthCheck": true,
		"debugPprof":  true,
	}
	s := NewWithConfig(config)
	if s.Configs()["port"] != config["port"] {
		t.Errorf("Port is different of default port %d != %d", s.Configs()["port"], config["port"])
	}
	if s.Configs()["host"] != config["host"] {
		t.Errorf("Host is different of default host %s != %s", s.Configs()["host"], config["host"])
	}
	if s.Configs()["basePath"] != config["basePath"] {
		t.Errorf("BasePath is different of default basePath %s != %s", s.Configs()["basePath"], config["basePath"])
	}
	if s.Configs()["stats"] != config["stats"] {
		t.Errorf("Stats is different of default stats %t != %t", s.Configs()["stats"], config["stats"])
	}
	if s.Configs()["healthCheck"] != config["healthCheck"] {
		t.Errorf("HealthCheck is different of default healthCheck %t != %t", s.Configs()["healthCheck"], config["healthCheck"])
	}
	if s.Configs()["debugPprof"] != config["debugPprof"] {
		t.Errorf("DebugPprof is different of default debugPprof %t != %t", s.Configs()["debugPprof"], config["debugPprof"])
	}

	config = map[string]interface{}{
		"stats":       true,
		"healthCheck": true,
		"debugPprof":  true,
	}

	s = NewWithConfig(config)
	if s.Configs()["port"] != defaultConfig["port"] {
		t.Errorf("Port is different of default port %d != %d", s.Configs()["port"], defaultConfig["port"])
	}
	if s.Configs()["host"] != defaultConfig["host"] {
		t.Errorf("Host is different of default host %s != %s", s.Configs()["host"], defaultConfig["host"])
	}
	if s.Configs()["basePath"] != defaultConfig["basePath"] {
		t.Errorf("BasePath is different of default basePath %s != %s", s.Configs()["basePath"], defaultConfig["basePath"])
	}
	if s.Configs()["stats"] != config["stats"] {
		t.Errorf("Stats is different of default stats %t != %t", s.Configs()["stats"], config["stats"])
	}
	if s.Configs()["healthCheck"] != config["healthCheck"] {
		t.Errorf("HealthCheck is different of default healthCheck %t != %t", s.Configs()["healthCheck"], config["healthCheck"])
	}
	if s.Configs()["debugPprof"] != config["debugPprof"] {
		t.Errorf("DebugPprof is different of default debugPprof %t != %t", s.Configs()["debugPprof"], config["debugPprof"])
	}
}

func Test_server_HandleError(t *testing.T) {
	f := func(w *worker.Worker, err error) {}
	s := New().HandleError(f)
	if s.handleError != nil && fmt.Sprintf("%p", f) != fmt.Sprintf("%p", s.handleError) {
		t.Errorf("HandleError is different of f %v != %v", fmt.Sprintf("%p", f), fmt.Sprintf("%p", s.handleError))
	}
}

func Test_server_HealthCheck(t *testing.T) {
	s := New().HealthCheck()
	if !s.Configs()["healthCheck"].(bool) {
		t.Error("HealthCheck setup on WorkerServer and his not enable")
	}
}

func Test_server_Stats(t *testing.T) {
	s := New().Stats()
	if !s.Configs()["stats"].(bool) {
		t.Error("Stats setup on WorkerServer and his not enable")
	}
}

func Test_server_DebugPprof(t *testing.T) {
	s := New().DebugPprof()
	if !s.Configs()["debugPprof"].(bool) {
		t.Error("DebugPprof setup on WorkerServer and his not enable")
	}
}

func Test_server_Run(t *testing.T) {
	s := New().Worker("w1", func(ctx context.Context) error { return nil })
	go func() {
		if err := s.Run(); err != nil {
			t.Errorf("Error when run WorkerServer %v", err)
		}
	}()
	s = New().HealthCheck().DebugPprof().Stats().Worker("w2", func(ctx context.Context) error { return nil })
	go func() {
		if err := s.Run(); err != nil {
			t.Errorf("Error when run WorkerServer %v", err)
		}
	}()

	<-time.After(time.Second)
	if err := s.GracefulStop(); err != nil {
		t.Error(err)
	}
}

func Test_server_Run_Error(t *testing.T) {
	if err := New().Run(); err != nil {
		t.Errorf("Error when run WorkerServer %v", err)
	}
}

func Test_server_Worker(t *testing.T) {
	nameWorker := "w1"
	handleWorker := func(ctx context.Context) error {
		<-time.After(3 * time.Second)
		return nil
	}
	concurrencyWorker := 1
	s := New().Worker(nameWorker, handleWorker)
	workers := s.Workers()
	if len(workers) != 1 {
		t.Error("Number of workers is different from setup")
	}
	w := workers[0]
	if w.Name != nameWorker {
		t.Errorf("Name of worker if different from setup %s != %s", w.Name, nameWorker)
	}
	if w.Handle != nil && fmt.Sprintf("%p", w.Handle) != fmt.Sprintf("%p", handleWorker) {
		t.Errorf("Name of worker if different from setup %p != %p", w.Handle, handleWorker)
	}
	if w.Concurrency != concurrencyWorker {
		t.Errorf("Concurrency of worker if different from setup %d != %d", w.Concurrency, concurrencyWorker)
	}
	if w.RestartAlways != false {
		t.Errorf("RestartAlaways of worker if different from setup %t != false", w.RestartAlways)
	}
	go func() {
		if err := s.Run(); err != nil {
			t.Error(err)
		}
	}()
	<-time.After(1 * time.Second)
	if !s.Healthy() {
		t.Error("Healthy expected is true but returned false")
	}
}

func Test_workerServer_StatsFunc(t *testing.T) {
	s := New().StatsFunc(func(writer http.ResponseWriter, request *http.Request) {})
	if _, ok := s.config["statsFunc"]; !ok {
		t.Error("StatsFunc is setup but still nil")
	}
}

func Test_workerServer_HealthCheckFunc(t *testing.T) {
	s := New().HealthCheckFunc(func(writer http.ResponseWriter, request *http.Request) {})
	if _, ok := s.config["healthCheckFunc"]; !ok {
		t.Error("HealthCheckFunc is setup but still nil")
	}
}

func TestWorkerServer_Healthy(t *testing.T) {
	s := New().CheckHealth(func() bool {
		return true
	})
	if !s.Healthy() {
		t.Error("Healthy expected is true but returned false")
	}
	s = New().CheckHealth(func() bool {
		return false
	})
	if s.Healthy() {
		t.Error("Healthy expected is false but returned true")
	}
	s = New()
	if !s.Healthy() {
		t.Error("Healthy expected is true but returned false")
	}
}

func TestWorkerServer_Infos(t *testing.T) {
	if infos := New().Infos(); infos == nil {
		t.Error("Infos expected is but returned nil")
	}
}

func TestWorkerServer_RunFullFeatures(t *testing.T) {
	err := New().Worker("test",
		func(ctx context.Context) error {
			<-time.After(3 * time.Second)
			return nil
		},
		worker.WithTimeout(2*time.Second),
		worker.WithDeadline(time.Now().Add(4*time.Second)),
		worker.WithConcurrency(2)).Run()
	if err != nil {
		t.Error(err)
	}
}

func TestExample_simpleWorker(t *testing.T) {
	server := New().
		Stats().
		HealthCheck().
		DebugPprof().
		HandleError(func(w *worker.Worker, err error) {
			log.Printf("Worker [%s] error: %s", w.Name, err)
		}).
		Worker(
			"w1",
			func(ctx context.Context) error {
				select {
				case <-time.After(5 * time.Second):
					return errors.New("test w1")
				case <-ctx.Done():
					return nil
				}
			},
			worker.WithRestartAlways()).
		Worker(
			"w2",
			func(ctx context.Context) error {
				select {
				case <-time.After(10 * time.Second):
					return nil
				case <-ctx.Done():
					return nil
				}
			}).
		Worker(
			"w3",
			func(ctx context.Context) error {
				select {
				case <-time.After(15 * time.Second):
					return errors.New("test w3")
				case <-ctx.Done():
					return nil
				}
			})
	go func() {
		if err := server.Run(); err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(16 * time.Second)
	if err := server.GracefulStop(); err != nil {
		t.Error(err)
	}
}
