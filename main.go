package main

import (
  "bufio"
  "context"
  "fmt"
  "log"
  "os"
  "io/ioutil"
  "strings"
  "net/http"
  "path/filepath"

  "voice-notetaking-app/pkg/database/sqlite"
  "voice-notetaking-app/service/speechtotext"
  "voice-notetaking-app/service/insight"
  "voice-notetaking-app/service/summarization"
  "voice-notetaking-app/service/tagging"


  
  "github.com/sashabaranov/go-openai"
)



// Counter variables for generating unique IDs
var (
  nodeIDCounter   int64 = 1
  edgeIDCounter   int64 = 1
  vertexIDCounter int64 = 1
)

// Define the KnowledgeGraph struct
type KnowledgeGraph struct {
  Nodes    map[int64]*Node
  Edges    map[int64]*Edge
  Vertices map[int64]*Vertex
}

// Define the Node struct
type Node struct {
  ID       int64
  Text     string
  Concepts []string
}

// Define the Edge struct
type Edge struct {
  SourceID int64
  TargetID int64
  Weight   float64
}

// Define the Vertex struct
type Vertex struct {
  ID       int64
  NodeID   int64
  TargetID int64
  Concept  string
}

// Define the Graph struct
type Graph struct {
  Nodes    []Node
  Edges    []Edge
  Vertices []Vertex
}


// Helper functions for generating unique IDs
func generateNodeID() int64 {
  nodeIDCounter++
  return nodeIDCounter
}

func generateEdgeID() int64 {
  edgeIDCounter++
  return edgeIDCounter
}

func generateVertexID() int64 {
  vertexIDCounter++
  return vertexIDCounter
}


// CalculateWeight calculates the weight between two sets of concepts based on Jaccard similarity
func calculateWeight(concepts1, concepts2 []string) float64 {
    // Convert concept slices to sets for easier comparison
    set1 := make(map[string]struct{})
    set2 := make(map[string]struct{})

    for _, concept := range concepts1 {
        set1[concept] = struct{}{}
    }

    for _, concept := range concepts2 {
        set2[concept] = struct{}{}
    }

    // Calculate Jaccard similarity
    intersection := 0
    for concept := range set1 {
        if _, exists := set2[concept]; exists {
            intersection++
        }
    }

    union := len(set1) + len(set2) - intersection

    // Prevent division by zero
    if union == 0 {
        return 0.0
    }

    return float64(intersection) / float64(union)
}

// Contains checks if a string exists in a slice
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

// getAllConcepts retrieves all unique concepts from the graph
func getAllConcepts(graph *Graph) []string {
    uniqueConcepts := make(map[string]struct{})

    for _, node := range graph.Nodes {
        for _, concept := range node.Concepts {
            uniqueConcepts[concept] = struct{}{}
        }
    }

    concepts := make([]string, 0, len(uniqueConcepts))
    for concept := range uniqueConcepts {
        concepts = append(concepts, concept)
    }

    return concepts
}




// ExtractNodesEdgesVertices uses AI to extract nodes, edges, and vertices based on the provided concepts
func ExtractNodesEdgesVertices(graph *Graph, client *openai.Client) error {
    // Prepare system prompt
    prompt := "You are an AI assistant that is an expert at understanding context. You will take the provided concepts and extract various nodes, edges, and vertices for our knowledge graph based on sentiment. Our aim is to allow the nodes, edges, and vertices to be parsed for added context, so keep that in mind. Respond with the nodes and any potential edges or vertices you want to add to the knowledge graph. Do not respond with the words edge, nodes, vertices, just the words themselves."

    // Get all concepts as a single string
    conceptsString := strings.Join(getAllConcepts(graph), ", ")

    // Extract nodes, edges, and vertices
    resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
        Model: openai.GPT3Dot5Turbo,
        Messages: []openai.ChatCompletionMessage{
            {
                Role:    openai.ChatMessageRoleUser,
                Content: prompt + "\nConcepts: " + conceptsString,
            },
        },
    })
    if err != nil {
        return fmt.Errorf("failed to extract nodes, edges, and vertices: %v", err)
    }

    // Check if messages are nil
    if resp.Messages == nil {
        return fmt.Errorf("no messages received in response")
    }

    // Process response and update the graph
    for _, message := range resp.Messages {
        if message.Role == openai.ChatMessageRoleAI {
            messageContent := message.Content
            // Then parse the message content and update the graph accordingly
        }
    }

    return nil
}






// NewKnowledgeGraph creates a new instance of KnowledgeGraph
func NewKnowledgeGraph() *KnowledgeGraph {
  return &KnowledgeGraph{
    Nodes:    make(map[int64]*Node),
    Edges:    make(map[int64]*Edge),
    Vertices: make(map[int64]*Vertex),
  }
}

func main() {
  // Initialize SQLite database
  dbPath := "./db/mydatabase.db"
  err := sqlite.Initialize(dbPath)
  if err != nil {
    log.Fatalf("Failed to initialize SQLite database: %v", err)
  }
  defer sqlite.Close()
  log.Println("SQLite database initialized successfully")

  // Read audio file
  audioFilePath := "testaudio.mp3" // Path to the audio file
  audioBytes, err := ioutil.ReadFile(audioFilePath)
  if err != nil {
    log.Fatalf("Failed to read audio file: %v", err)
  }
  log.Println("Audio file read successfully")

  // Convert audio to text using speech-to-text service
  transcription, err := speechtotext.UploadToAssemblyAI(audioBytes)
  if err != nil {
    log.Printf("Failed to transcribe audio: %v", err)
  }
  log.Println("Transcription:", transcription)

  // Load knowledge graph data
  graph, err := LoadGraph("knowledge_graph.txt")
  if err != nil {
    log.Fatalf("Failed to load knowledge graph: %v", err)
  }
  log.Println("Knowledge graph loaded successfully")

  // HTTP handler to upload voice note
  http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
    // Parse multipart form data
    err := r.ParseMultipartForm(10 << 20) // 10 MB
    if err != nil {
      log.Printf("Failed to parse form data: %v", err)
      http.Error(w, "Failed to parse form data", http.StatusBadRequest)
      return
    }
    log.Println("Form data parsed successfully")

    // Get audio file from form data
    file, _, err := r.FormFile("audio")
    if err != nil {
      log.Printf("Failed to read audio file: %v", err)
      http.Error(w, "Failed to read audio file", http.StatusBadRequest)
      return
    }
    defer file.Close()
    log.Println("Audio file read successfully")

    // Read file content
    fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
      log.Printf("Failed to read file content: %v", err)
      http.Error(w, "Failed to read file content", http.StatusInternalServerError)
      return
    }
    log.Println("File content read successfully")

    // Convert audio to text using speech-to-text service
    transcription, err := speechtotext.UploadToAssemblyAI(fileBytes)
    if err != nil {
      log.Printf("Failed to transcribe audio: %v", err)
      http.Error(w, "Failed to transcribe audio", http.StatusInternalServerError)
      return
    }
    log.Println("Transcription:", transcription)

    // Summarize the transcription
    summary, err := summarization.SummarizeText(transcription)
    if err != nil {
      log.Printf("Failed to summarize text: %v", err)
      http.Error(w, "Failed to summarize text", http.StatusInternalServerError)
      return
    }
    log.Println("Summary:", summary)

    // Tag the transcription
    tags, err := tagging.TagText(transcription)
    if err != nil {
      log.Printf("Failed to tag text: %v", err)
      http.Error(w, "Failed to tag text", http.StatusInternalServerError)
      return
    }
    log.Println("Tags:", tags)

    // Generate insights from the transcription
    groupedNotes := []string{transcription} // For simplicity, using the entire transcription as one note
    insightText, err := insight.GenerateInsight(groupedNotes)
    if err != nil {
      log.Printf("Failed to generate insight: %v", err)
      http.Error(w, "Failed to generate insight", http.StatusInternalServerError)
      return
    }
    log.Println("Insight:", insightText)

    // Insert the transcription into the database
    userID := 1                                  // Example user ID
    filePath := filepath.Join("recordings", "audio.mp3") // Example file path
    recordingID, err := sqlite.InsertRecording(int64(userID), filePath, transcription)
    if err != nil {
      log.Printf("Failed to insert recording into database: %v", err)
      http.Error(w, "Failed to insert recording into database", http.StatusInternalServerError)
      return
    }
    log.Println("Recording inserted with ID:", recordingID)

    // Build or update knowledge graph with the provided note text and concepts
    if err := BuildOrUpdateKnowledgeGraph(&graph, transcription, tags); err != nil {
      log.Fatalf("Failed to build or update knowledge graph: %v", err)
    }
    log.Println("Knowledge graph updated successfully")

    // Save the updated graph
    if err := SaveGraph(&graph, "knowledge_graph.txt"); err != nil {
      log.Fatalf("Failed to save knowledge graph: %v", err)
    }
    log.Println("Knowledge graph saved successfully")

    // Send success response
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Processing completed successfully!")
  })

  // Start HTTP server
  log.Println("Server is running on port 8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}

// BuildOrUpdateKnowledgeGraph builds the knowledge graph with the provided note text and tags, or updates an existing graph
func BuildOrUpdateKnowledgeGraph(graph *Graph, noteText string, tags []string) error {
  // Extract concepts from tags (assuming tags represent concepts)
  concepts := tags

  // Create nodes for the note
  node := Node{
    ID:       generateNodeID(),
    Text:     noteText,
    Concepts: concepts,
  }
  graph.Nodes = append(graph.Nodes, node)

  // Create edges and vertices based on the relationships between nodes
  for _, existingNode := range graph.Nodes {
    if existingNode.ID != node.ID {
      // Calculate edge weight based on concept similarity
      weight := calculateWeight(node.Concepts, existingNode.Concepts)
      if weight > 0 {
        // Create an edge between the nodes
        edge := Edge{
          SourceID: node.ID,
          TargetID: existingNode.ID,
          Weight:   weight,
        }
        graph.Edges = append(graph.Edges, edge)

        // Create vertices for the concepts shared by the nodes
        for _, concept := range node.Concepts {
          if contains(existingNode.Concepts, concept) {
            vertex := Vertex{
              NodeID:   node.ID,
              TargetID: existingNode.ID,
              Concept:  concept,
            }
            graph.Vertices = append(graph.Vertices, vertex)
          }
        }
      }
    }
  }

  return nil
}


// SaveGraph saves the knowledge graph to a file.
func SaveGraph(graph *Graph, filename string) error {
  filePath := filename

  // Open the file for writing
  file, err := os.Create(filePath)
  if err != nil {
    return fmt.Errorf("failed to create file: %v", err)
  }
  defer file.Close()

  // Write nodes to the file
  for _, node := range graph.Nodes {
    _, err := fmt.Fprintf(file, "Node ID: %d\nText: %s\nConcepts: %s\n\n", node.ID, node.Text, strings.Join(node.Concepts, ", "))
    if err != nil {
      return fmt.Errorf("failed to write node to file: %v", err)
    }
  }

  // Write edges to the file
  for _, edge := range graph.Edges {
    _, err := fmt.Fprintf(file, "Edge: SourceID: %d, TargetID: %d, Weight: %f\n", edge.SourceID, edge.TargetID, edge.Weight)
    if err != nil {
      return fmt.Errorf("failed to write edge to file: %v", err)
    }
  }

  return nil
}

// LoadGraph loads the knowledge graph from the database
func LoadGraph(filename string) (Graph, error) {
  filePath := filename
  var graph Graph

  // Open the file for reading
  file, err := os.Open(filePath)
  if err != nil {
    return graph, fmt.Errorf("failed to open file: %v", err)
  }
  defer file.Close()

  var node Node
  var edge Edge

  // Read lines from the file and construct the graph
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    if strings.HasPrefix(line, "Node ID:") {
      if node.ID != 0 {
        graph.Nodes = append(graph.Nodes, node)
        node = Node{}
      }
      _, err := fmt.Sscanf(line, "Node ID: %d", &node.ID)
      if err != nil {
        return graph, fmt.Errorf("failed to parse node ID: %v", err)
      }
    } else if strings.HasPrefix(line, "Text:") {
      _, err := fmt.Sscanf(line, "Text: %s", &node.Text)
      if err != nil {
        return graph, fmt.Errorf("failed to parse node text: %v", err)
      }
    } else if strings.HasPrefix(line, "Concepts:") {
      conceptsStr := strings.TrimPrefix(line, "Concepts: ")
      node.Concepts = strings.Split(conceptsStr, ", ")
    } else if strings.HasPrefix(line, "Edge:") {
      _, err := fmt.Sscanf(line, "Edge: SourceID: %d, TargetID: %d, Weight: %f", &edge.SourceID, &edge.TargetID, &edge.Weight)
      if err != nil {
        return graph, fmt.Errorf("failed to parse edge: %v", err)
      }
      graph.Edges = append(graph.Edges, edge)
    }
  }

  if err := scanner.Err(); err != nil {
    return graph, fmt.Errorf("error reading file: %v", err)
  }

  return graph, nil
}
