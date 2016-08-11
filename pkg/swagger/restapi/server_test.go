package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/workflow-manager-api/config"
	"github.com/deis/workflow-manager-api/pkg/data"
	"github.com/deis/workflow-manager-api/pkg/handlers"
	"github.com/deis/workflow-manager-api/pkg/swagger/models"
	"github.com/deis/workflow-manager-api/pkg/swagger/restapi/operations"
	"github.com/deis/workflow-manager-api/rest"
	spec "github.com/go-swagger/go-swagger/spec"
	"github.com/go-swagger/go-swagger/swag"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
)

const (
	componentName     = "testcomponent"
	clusterID         = "testcluster"
	releaseTimeFormat = "2006-01-02T15:04:05Z"
)

var (
	nowTime          = time.Now()
	futureTime       = nowTime.Add(1 * time.Hour)
	pastTime         = nowTime.Add(-1 * time.Hour)
	doctorReportUUID = uuid.New()
)

func newServer(db *gorm.DB) (*httptest.Server, error) {
	swaggerSpec, err := spec.New(SwaggerJSON, "")
	if err != nil {
		return nil, err
	}
	api := operations.NewWorkflowManagerAPI(swaggerSpec)
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		swag.CommandLineOptionsGroup{
			LongDescription:  "deisUnitTests",
			Options:          GormDb{db},
			ShortDescription: "deisUnitTests",
		},
	}
	// Routes consist of a path and a handler function.
	return httptest.NewServer(configureAPI(api)), nil
}

func urlPath(ver string, remainder ...string) string {
	return fmt.Sprintf("%s/%s", ver, strings.Join(remainder, "/"))
}

// tests the GET /{apiVersion}/versions/{train}/{component}/{version} endpoint
func TestGetVersion(t *testing.T) {
	db, err := data.NewMemDB()
	assert.NoErr(t, err)
	assert.NoErr(t, data.VerifyPersistentStorage(db))
	srv, err := newServer(db)
	assert.NoErr(t, err)
	defer srv.Close()
	componentVer := models.ComponentVersion{
		Component: &models.Component{Name: componentName},
		Version: &models.Version{Train: "beta", Version: "2.0.0-beta-2", Released: "2016-03-31T23:54:39Z", Data: &models.VersionData{
			Description: "release notes",
		}},
	}
	_, err = data.UpsertVersion(db, componentVer)
	assert.NoErr(t, err)
	resp, err := httpGet(srv, urlPath("v3", "versions", componentVer.Version.Train, componentVer.Component.Name, componentVer.Version.Version))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
	decodedVer := new(models.ComponentVersion)
	assert.NoErr(t, json.NewDecoder(resp.Body).Decode(decodedVer))
	assert.Equal(t, *decodedVer, componentVer, "component version")
}

// tests the GET /{apiVersion}/versions/{train}/{component} endpoint
func TestGetComponentTrainVersions(t *testing.T) {
	memDB, err := data.NewMemDB()
	assert.NoErr(t, err)
	assert.NoErr(t, data.VerifyPersistentStorage(memDB))
	srv, err := newServer(memDB)
	assert.NoErr(t, err)
	defer srv.Close()
	componentVers := []models.ComponentVersion{}
	componentVer1 := models.ComponentVersion{
		Component: &models.Component{Name: componentName},
		Version: &models.Version{Train: "beta", Version: "2.0.0-beta-1", Released: "2016-03-30T23:54:39Z", Data: &models.VersionData{
			Description: "release notes",
		}},
	}
	componentVer2 := models.ComponentVersion{
		Component: &models.Component{Name: componentName},
		Version: &models.Version{Train: "beta", Version: "2.0.0-beta-2", Released: "2016-03-31T23:54:39Z", Data: &models.VersionData{
			Description: "release notes",
		}},
	}
	componentVers = append(componentVers, componentVer1)
	componentVers = append(componentVers, componentVer2)
	_, err = data.UpsertVersion(memDB, componentVers[0])
	assert.NoErr(t, err)
	_, err = data.UpsertVersion(memDB, componentVers[1])
	assert.NoErr(t, err)
	resp, err := httpGet(srv, urlPath("v3", "versions", componentVer1.Version.Train, componentVer1.Component.Name))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
	decodedVer := new(operations.GetComponentByNameOKBodyBody)
	assert.NoErr(t, json.NewDecoder(resp.Body).Decode(decodedVer))
	assert.Equal(t, *decodedVer.Data[0], componentVers[0], "component versions")
	assert.Equal(t, *decodedVer.Data[1], componentVers[1], "component versions")
}

// tests the GET /{apiVersion}/versions/{train}/{component}/latest endpoint
func TestGetLatestComponentTrainVersion(t *testing.T) {
	const componentName = "testcomponent"
	const train = "testtrain"
	memDB, err := data.NewMemDB()
	assert.NoErr(t, err)
	assert.NoErr(t, data.VerifyPersistentStorage(memDB))
	srv, err := newServer(memDB)
	assert.NoErr(t, err)
	defer srv.Close()

	const numCVs = 4
	const latestCVIdx = 2
	componentVersions := make([]models.ComponentVersion, numCVs)
	for i := 0; i < numCVs; i++ {
		cv := models.ComponentVersion{}
		cv.Component = &models.Component{}
		cv.Version = &models.Version{}
		cv.Component.Name = componentName
		desc := fmt.Sprintf("description%d", i)
		cv.Component.Description = &desc
		cv.Version.Train = train
		cv.Version.Version = fmt.Sprintf("testversion%d", i)
		cv.Version.Released = time.Now().Add(time.Duration(i) * time.Hour).Format(releaseTimeFormat)
		cv.Version.Data = &models.VersionData{
			Description: fmt.Sprintf("data%d", i),
		}
		if i == latestCVIdx {
			cv.Version.Released = time.Now().Add(time.Duration(numCVs+1) * time.Hour).Format(releaseTimeFormat)
		}
		if _, setErr := data.UpsertVersion(memDB, cv); setErr != nil {
			t.Fatalf("error setting component version %d (%s)", i, setErr)
		}
		componentVersions[i] = cv
	}
	path := urlPath("v3", "versions", train, componentName, "latest")
	resp, err := httpGet(srv, path)
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
	cv := new(models.ComponentVersion)
	assert.NoErr(t, json.NewDecoder(resp.Body).Decode(cv))
	exCV := componentVersions[latestCVIdx]

	assert.Equal(t, cv.Component.Name, exCV.Component.Name, "component name")
	// since the versions table doesn't store a description now, make sure it comes back empty
	assert.Nil(t, cv.Component.Description, "component name")

	assert.Equal(t, cv.Version.Train, exCV.Version.Train, "component version")
	assert.Equal(t, cv.Version.Version, exCV.Version.Version, "component version")
	assert.Equal(t, cv.Version.Released, exCV.Version.Released, "component release time")
	assert.Equal(t, cv.Version.Data, exCV.Version.Data, "component version data")
}

// tests the POST /{apiVersion}/versions/{train}/{component}/{version} endpoint
func TestPostVersions(t *testing.T) {
	memDB, err := data.NewMemDB()
	assert.NoErr(t, err)
	assert.NoErr(t, data.VerifyPersistentStorage(memDB))
	srv, err := newServer(memDB)
	assert.NoErr(t, err)
	defer srv.Close()
	train := "beta"
	version := "2.0.0-beta-2"
	componentVer := models.ComponentVersion{
		Component: &models.Component{Name: componentName},
		Version: &models.Version{Train: train, Version: version, Released: "2016-03-31T23:54:39Z", Data: &models.VersionData{
			Description: "release notes",
		}},
	}
	body := new(bytes.Buffer)
	assert.NoErr(t, json.NewEncoder(body).Encode(componentVer))
	resp, err := httpPost(srv, urlPath("v3", "versions", train, componentName, version), string(body.Bytes()))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
	retComponentVersion := new(models.ComponentVersion)
	assert.NoErr(t, json.NewDecoder(resp.Body).Decode(retComponentVersion))
	// TODO: version data property not traveling and returning as expected
	assert.Equal(t, *retComponentVersion, componentVer, "component version")
	fetchedComponentVersion, err := data.GetVersion(memDB, componentVer)
	assert.NoErr(t, err)
	assert.Equal(t, fetchedComponentVersion, componentVer, "component version")
}

// tests the GET /{apiVersion}/clusters/count endpoint
func TestGetClusters(t *testing.T) {
	memDB, err := data.NewMemDB()
	assert.NoErr(t, err)
	assert.NoErr(t, data.VerifyPersistentStorage(memDB))
	srv, err := newServer(memDB)
	assert.NoErr(t, err)
	defer srv.Close()
	resp, err := httpGet(srv, urlPath("v3", "clusters", "count"))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
}

// tests the GET /{apiVersion}/clusters/{id} endpoint
func TestGetClusterByID(t *testing.T) {
	memDB, err := data.NewMemDB()
	assert.NoErr(t, err)
	assert.NoErr(t, data.VerifyPersistentStorage(memDB))
	srv, err := newServer(memDB)
	assert.NoErr(t, err)
	defer srv.Close()
	cluster := models.Cluster{}
	cluster.ID = clusterID
	cluster.Components = nil
	newCluster, err := data.UpsertCluster(memDB, clusterID, cluster)
	assert.NoErr(t, err)
	resp, err := httpGet(srv, urlPath("v3", "clusters", clusterID))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, 200, "response code")
	decodedCluster := new(models.Cluster)
	assert.NoErr(t, json.NewDecoder(resp.Body).Decode(decodedCluster))
	assert.Equal(t, *decodedCluster, newCluster, "returned cluster")
}

// tests the POST {apiVersion}/clusters/{id} endpoint
func TestPostClusters(t *testing.T) {
	jsonData := `{"Components": [{"Component": {"Name": "component-a"}, "Version": {"Version": "1.0"}}], "ID": "testcluster"}`
	memDB, err := data.NewMemDB()
	assert.NoErr(t, err)
	assert.NoErr(t, data.VerifyPersistentStorage(memDB))
	srv, err := newServer(memDB)
	assert.NoErr(t, err)
	defer srv.Close()
	resp, err := httpPost(srv, urlPath("v3", "clusters"), jsonData)
	if err != nil {
		t.Fatalf("POSTing to endpoint (%s)", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d", resp.StatusCode)
	}
	resp, err = httpGet(srv, urlPath("v3", "clusters", clusterID))
	defer resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d", resp.StatusCode)
	}
	cluster := new(models.Cluster)
	if err := json.NewDecoder(resp.Body).Decode(cluster); err != nil {
		t.Fatalf("error reading response body (%s)", err)
	}
	if len(cluster.Components) <= 0 {
		t.Fatalf("no components returned")
	}
	if cluster.Components[0].Component.Name != "component-a" {
		t.Error("unexpected component name from JSON response")
	}
	// Note that we have to dereference "Version" twice because cluster.Components[0].Version
	// is itself a models.Version, which has both a "Released" and "Version" field
	if cluster.Components[0].Version.Version != "1.0" {
		t.Error("unexpected component version from JSON response")
	}
}

func TestGetLatestVersions(t *testing.T) {
	memDB, err := data.NewMemDB()
	assert.NoErr(t, err)
	assert.NoErr(t, data.VerifyPersistentStorage(memDB))
	assert.NoErr(t, err)
	srv, err := newServer(memDB)
	assert.NoErr(t, err)
	defer srv.Close()
	const numComponentVersions = 6
	components := make(map[string]models.ComponentVersion)
	for i := 0; i < numComponentVersions; i++ {
		name := fmt.Sprintf("component%d", i)
		train := fmt.Sprintf("train%d", i)
		releaseTime1 := time.Now().Add(time.Duration(i+1) * time.Hour)
		releaseTime2 := time.Now().Add(time.Duration((i+1)*2) * time.Hour)
		cv1 := models.ComponentVersion{
			Component: &models.Component{Name: name},
			Version: &models.Version{
				Train:    train,
				Released: releaseTime1.Format(releaseTimeFormat),
				Version:  fmt.Sprintf("version%d-1", i),
			},
		}
		cv2 := models.ComponentVersion{
			Component: &models.Component{Name: name},
			Version: &models.Version{
				Train:    train,
				Released: releaseTime2.Format(releaseTimeFormat),
				Version:  fmt.Sprintf("version%d-2", i),
			},
		}

		if _, err := data.UpsertVersion(memDB, cv1); err != nil {
			t.Fatalf("Error setting component %d (%s)", i, err)
		}
		if _, err := data.UpsertVersion(memDB, cv2); err != nil {
			t.Fatalf("Error setting component %d (%s)", i, err)
		}
		components[cv2.Component.Name] = cv2
	}

	postBody := handlers.SparseComponentAndTrainInfoJSONWrapper{
		Data: make([]handlers.SparseComponentAndTrainInfo, len(components)),
	}
	i := 0
	for _, component := range components {
		postBody.Data[i] = handlers.SparseComponentAndTrainInfo{
			Component: handlers.SparseComponentInfo{Name: component.Component.Name},
			Version:   handlers.SparseVersionInfo{Train: component.Version.Train},
		}
		i++
	}
	postBodyReader := new(bytes.Buffer)
	assert.NoErr(t, json.NewEncoder(postBodyReader).Encode(postBody))

	resp, err := httpPost(srv, urlPath("v3", "versions", "latest"), string(postBodyReader.Bytes()))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
	respStruct := new(handlers.ComponentVersionsJSONWrapper)
	assert.NoErr(t, json.NewDecoder(resp.Body).Decode(respStruct))
	assert.Equal(t, len(respStruct.Data), len(components), "length of response")
	for _, component := range respStruct.Data {
		expected, ok := components[component.Component.Name]
		assert.True(t, ok, "%s not found in the component/train response body", component.Component.Name)
		assert.Equal(t, component.Version.Version, expected.Version.Version, "component version")
		assert.Equal(t, component.Version.Train, expected.Version.Train, "component train")
		assert.Equal(t, component.Version.Released, expected.Version.Released, "component released time")
	}
}

func TestFilterByClusterAge(t *testing.T) {
	filter := data.ClusterAgeFilter{
		CheckedInBefore: timeFuture().Add(2 * time.Hour),
		CheckedInAfter:  timePast(),
		CreatedAfter:    timePast().Add(-1 * time.Hour),
		CreatedBefore:   timeFuture(),
	}
	memDB, err := data.NewMemDB()
	assert.NoErr(t, err)
	assert.NoErr(t, data.VerifyPersistentStorage(memDB))
	cluster := models.Cluster{}
	cluster.ID = uuid.New()
	srv, err := newServer(memDB)
	assert.NoErr(t, err)
	defer srv.Close()
	_, setErr := data.UpsertCluster(memDB, cluster.ID, cluster)
	assert.NoErr(t, setErr)
	assert.NoErr(t, data.CheckInCluster(memDB, cluster.ID, time.Now(), cluster))
	queryPairsMap := map[string]string{
		rest.CheckedInBeforeQueryStringKey: filter.CheckedInBefore.Format(data.StdTimestampFmt),
		rest.CheckedInAfterQueryStringKey:  filter.CheckedInAfter.Format(data.StdTimestampFmt),
		rest.CreatedBeforeQueryStringKey:   filter.CreatedBefore.Format(data.StdTimestampFmt),
		rest.CreatedAfterQueryStringKey:    filter.CreatedAfter.Format(data.StdTimestampFmt),
	}
	queryPairs := make([]string, len(queryPairsMap))
	i := 0
	for k, v := range queryPairsMap {
		queryPairs[i] = fmt.Sprintf("%s=%s", k, v)
		i++
	}

	route := urlPath("v3", "clusters", "age")
	route += fmt.Sprintf("?%s", strings.Join(queryPairs, "&"))
	resp, err := httpGet(srv, route)
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
	var respEnvelope struct {
		Data []models.Cluster `json:"data"`
	}
	assert.NoErr(t, json.NewDecoder(resp.Body).Decode(&respEnvelope))
	assert.Equal(t, len(respEnvelope.Data), 1, "length of the clusters list")
	assert.Equal(t, respEnvelope.Data[0].ID, cluster.ID, "returned cluster ID")
}

// tests the POST /{apiVersion}/doctor/{id} endpoint
func TestPostDoctor(t *testing.T) {
	memDB, err := data.NewMemDB()
	assert.NoErr(t, err)
	assert.NoErr(t, data.VerifyPersistentStorage(memDB))
	srv, err := newServer(memDB)
	assert.NoErr(t, err)
	defer srv.Close()
	jsonData := `{"workflow":{"components":[{"component":{"description":"Deis Workflow","name":"deis-builder"},"version":{"train":"","version":"v2.0.0"}},{"component":{"description":"Deis Workflow","name":"deis-controller"},"version":{"train":"","version":"v2.0.0"}},{"component":{"description":"Deis Workflow","name":"deis-database"},"version":{"train":"","version":"v2.0.0"}},{"component":{"description":"Deis Workflow","name":"deis-logger"},"version":{"train":"","version":"v2.0.0"}},{"component":{"description":"Deis Workflow","name":"deis-minio"},"version":{"train":"","version":"v2.0.0"}},{"component":{"description":"Deis Workflow","name":"deis-monitor-grafana"},"version":{"train":"","version":"v2.0.0"}},{"component":{"description":"Deis Workflow","name":"deis-monitor-influxdb"},"version":{"train":"","version":"v2.0.0"}},{"component":{"description":"Deis Workflow","name":"deis-monitor-stdout"},"version":{"train":"","version":"v2.0.0"}},{"component":{"description":"Deis Workflow","name":"deis-registry"},"version":{"train":"","version":"v2.0.0"}},{"component":{"description":"Deis Workflow","name":"deis-router"},"version":{"train":"","version":"v2.0.0"}},{"component":{"description":"Deis Workflow","name":"deis-workflow-manager"},"version":{"train":"","version":"v2.0.0"}}],"id":"6cd6539e-4225-43a1-89e7-0155b8ea1de6"}}`
	resp, err := httpPost(srv, urlPath("v3", "doctor", doctorReportUUID), jsonData)
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
}

// tests the GET /{apiVersion}/doctor/{id} endpoint
func TestGetDoctor(t *testing.T) {
	db, err := data.NewMemDB()
	assert.NoErr(t, err)
	assert.NoErr(t, data.VerifyPersistentStorage(db))
	srv, err := newServer(db)
	assert.NoErr(t, err)
	defer srv.Close()
	jsonData := fmt.Sprintf(`{"workflow":{"components":[{"component":{"description":"Deis Workflow","name":"deis-builder"},"version":{"train":"","version":"v2.0.0"}},{"component":{"description":"Deis Workflow","name":"deis-controller"},"version":{"train":"","version":"v2.0.0"}},{"component":{"description":"Deis Workflow","name":"deis-database"},"version":{"train":"","version":"v2.0.0"}},{"component":{"description":"Deis Workflow","name":"deis-logger"},"version":{"train":"","version":"v2.0.0"}},{"component":{"description":"Deis Workflow","name":"deis-minio"},"version":{"train":"","version":"v2.0.0"}},{"component":{"description":"Deis Workflow","name":"deis-monitor-grafana"},"version":{"train":"","version":"v2.0.0"}},{"component":{"description":"Deis Workflow","name":"deis-monitor-influxdb"},"version":{"train":"","version":"v2.0.0"}},{"component":{"description":"Deis Workflow","name":"deis-monitor-stdout"},"version":{"train":"","version":"v2.0.0"}},{"component":{"description":"Deis Workflow","name":"deis-registry"},"version":{"train":"","version":"v2.0.0"}},{"component":{"description":"Deis Workflow","name":"deis-router"},"version":{"train":"","version":"v2.0.0"}},{"component":{"description":"Deis Workflow","name":"deis-workflow-manager"},"version":{"train":"","version":"v2.0.0"}}],"id":"%s"}}`, clusterID)
	_, err = httpPost(srv, urlPath("v3", "doctor", doctorReportUUID), jsonData)
	assert.NoErr(t, err)
	resp, err := httpGetBasicAuth(srv, urlPath("v3", "doctor", doctorReportUUID), config.Spec.DBUser, config.Spec.DBPass)
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
	doctorInfoResponse := new(models.DoctorInfo)
	assert.NoErr(t, json.NewDecoder(resp.Body).Decode(doctorInfoResponse))
	assert.Equal(t, doctorInfoResponse.Workflow.ID, clusterID, "cluster ID")
}

func timeFuture() time.Time {
	return futureTime
}

func timePast() time.Time {
	return pastTime
}

func timeNow() time.Time {
	return nowTime
}

func httpGet(s *httptest.Server, route string) (*http.Response, error) {
	return http.Get(s.URL + "/" + route)
}

func httpPost(s *httptest.Server, route string, json string) (*http.Response, error) {
	fullURL := s.URL + "/" + route
	return http.Post(fullURL, "application/json", bytes.NewBuffer([]byte(json)))
}

func httpGetBasicAuth(s *httptest.Server, route, user, pass string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", s.URL+"/"+route, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(user, pass)
	return client.Do(req)
}

func httpPostBasicAuth(s *httptest.Server, route, json, user, pass string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", s.URL+"/"+route, bytes.NewBuffer([]byte(json)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(user, pass)
	return client.Do(req)
}
