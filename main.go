package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DkronJob represents a Dkron job.
type DkronJob struct {
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
	Timezone string `json:"timezone"`
	Owner    string `json:"owner"`
	Disable  bool   `json:"disable"`
	Tags     struct {
		Server string `json:"server"`
	} `json:"tags"`
	Metadata struct {
		User string `json:"user"`
	} `json:"metadata"`
	Concurrency      string `json:"concurrency"`
	Excecutor        string `json:"executor"`
	Excecutor_config struct {
		Command string `json:"command"`
	}
}

func main() {
	r := gin.Default()

	// Define endpoints
	r.GET("/jobs", listJobs)
	r.POST("/jobs", createJob)
	r.PUT("/jobs/:name", updateJob)
	r.DELETE("/jobs/:name", deleteJob)

	// Start server
	if err := r.Run(":8000"); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

func listJobs(c *gin.Context) {
	resp, err := http.Get("http://localhost:8080/v1/jobs")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Create a buffer to store the response body
	var body bytes.Buffer
	_, err = io.Copy(&body, resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Print the raw JSON response for debugging
	fmt.Println("Raw JSON response:", body.String())

	// Decode the JSON response into a slice of DkronJob objects
	var jobs []DkronJob
	if err := json.NewDecoder(&body).Decode(&jobs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

func createJob(c *gin.Context) {
	var job DkronJob
	if err := c.BindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jobData, err := json.Marshal(job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp, err := http.Post("http://localhost:8080/v1/jobs", "application/json", bytes.NewBuffer(jobData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	c.JSON(resp.StatusCode, gin.H{"message": "Job created"})
}

func updateJob(c *gin.Context) {
	name := c.Param("name")
	var job DkronJob
	if err := c.BindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jobData, err := json.Marshal(job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/v1/jobs/"+name, bytes.NewBuffer(jobData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	c.JSON(resp.StatusCode, gin.H{"message": "Job updated"})
}

func deleteJob(c *gin.Context) {
	name := c.Param("name")
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/v1/jobs/"+name, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	c.JSON(resp.StatusCode, gin.H{"message": "Job deleted"})
}
