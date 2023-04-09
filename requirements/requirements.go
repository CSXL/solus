package requirements

import (
	"fmt"
	"os"

	"github.com/CSXL/solus/ai"
	"github.com/CSXL/solus/ai/chat"
	"github.com/CSXL/solus/config"
	"github.com/joho/godotenv"
)

type RequirementsConfig struct {
	RequirementsPrompt string // The prompt to use when generating requirements
	OpenAIAPIKey       string // The OpenAI API key to use when generating requirements
}

func (r *RequirementsConfig) ToAIConfig() *ai.AIConfig {
	return ai.NewAIConfig(r.OpenAIAPIKey)
}

func NewRequirementsConfig(requirementsPrompt string, openAIAPIKey string) *RequirementsConfig {
	return &RequirementsConfig{
		RequirementsPrompt: requirementsPrompt,
		OpenAIAPIKey:       openAIAPIKey,
	}
}

func LoadRequirementsConfig() (*RequirementsConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	config_reader := config.New()
	err = config_reader.Read("requirements_config", ".")
	if err != nil {
		return nil, err
	}
	requirementsPrompt := config_reader.Get("requirements_prompt").(string)
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	requirements_config := NewRequirementsConfig(requirementsPrompt, openAIAPIKey)
	return requirements_config, nil
}

type RequirementsGenerator struct {
	Conversation          *chat.Conversation
	requirementsConfig    *RequirementsConfig
	inputConversation     string
	GeneratedRequirements string
}

// Creates a new RequirementsGenerator
// The requirements generator will use the given input data to generate the
// requirements for the project.
func NewRequirementsGenerator(inputConversation string, config *RequirementsConfig) *RequirementsGenerator {
	conversationName := "requirements"
	aiConfig := config.ToAIConfig()
	conversation := chat.NewConversation(conversationName, aiConfig)
	return &RequirementsGenerator{
		inputConversation:     inputConversation,
		Conversation:          conversation,
		requirementsConfig:    config,
		GeneratedRequirements: "",
	}
}

// Loads the conversation data from a file
func LoadInputConversation(filename string) (string, error) {
	fileContentBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	fileContentString := string(fileContentBytes)
	return fileContentString, nil
}

func (r *RequirementsGenerator) buildPrompt() string {
	inputPrompt := r.requirementsConfig.RequirementsPrompt
	beginConversation := "\n====BEGIN CONVERSATION====\n"
	inputConversation := r.inputConversation
	endConversation := "\n====END CONVERSATION====\n"
	prompt := fmt.Sprintf("%s%s%s%s", inputPrompt, beginConversation, inputConversation, endConversation)
	return prompt
}

func (r *RequirementsGenerator) promptModel() (string, error) {
	prompt := r.buildPrompt()
	r.Conversation.ResetMessages()
	responseMessage, err := r.Conversation.SendSystemMessage(prompt)
	if err != nil {
		return "", err
	}
	responseMessage.Serialize()
	responseContent := responseMessage.GetContent()
	return responseContent, nil
}

// Generates the requirements for the project
func (r *RequirementsGenerator) Generate() (string, error) {
	generatedRequirements, err := r.promptModel()
	if err != nil {
		return "", err
	}
	r.GeneratedRequirements = generatedRequirements
	return generatedRequirements, nil
}

// Saves the generated requirements to a file
func (r *RequirementsGenerator) SaveGeneratedRequirements(filename string) error {
	requirements := r.GeneratedRequirements
	err := os.WriteFile(filename, []byte(requirements), 0644)
	if err != nil {
		return err
	}
	return nil
}
