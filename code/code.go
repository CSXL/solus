package code

import (
	"fmt"
	"os"

	"github.com/CSXL/solus/ai"
	"github.com/CSXL/solus/ai/chat"
	"github.com/CSXL/solus/code/syncfiles"
	"github.com/CSXL/solus/config"
	"github.com/joho/godotenv"
)

type CodeConfig struct {
	CodePrompt       string // The prompt to use when generating code
	OpenAIAPIKey     string // The OpenAI API key to use when generating code
	GenerationFolder string // The folder to generate code in
}

func (c *CodeConfig) ToAIConfig() *ai.AIConfig {
	return ai.NewAIConfig(c.OpenAIAPIKey)
}

func NewCodeConfig(generationFolder string, openAIAPIKey string) *CodeConfig {
	return &CodeConfig{
		GenerationFolder: generationFolder,
		OpenAIAPIKey:     openAIAPIKey,
	}
}

func LoadCodeConfig() (*CodeConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	config_reader := config.New()
	err = config_reader.Read("code_config", ".")
	if err != nil {
		return nil, err
	}
	codePrompt := config_reader.Get("code_prompt").(string)
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	code_config := NewCodeConfig(codePrompt, openAIAPIKey)
	return code_config, nil
}

type CodeGenerator struct {
	Conversation *chat.Conversation
	codeConfig   *CodeConfig
	ProjectState string
}

// Creates a new CodeGenerator
// The code generator will use the given input data to generate the
// code for the project.
func NewCodeGenerator(generationFolder string, config *CodeConfig) *CodeGenerator {
	conversationName := "code"
	aiConfig := config.ToAIConfig()
	conversation := chat.NewConversation(conversationName, aiConfig)
	config.GenerationFolder = generationFolder
	return &CodeGenerator{
		Conversation: conversation,
		codeConfig:   config,
	}
}

func (c *CodeGenerator) buildPrompt() string {
	inputPrompt := c.codeConfig.CodePrompt
	beginCurrentState := "\n====CURRENT STATE====\n"
	currentState := c.ProjectState
	endCurrentState := "\n====END CURRENT STATE====\n"
	prompt := fmt.Sprintf("%s%s%s%s", inputPrompt, beginCurrentState, currentState, endCurrentState)
	return prompt
}

func (c *CodeGenerator) loadProjectState() error {
	projectState, err := syncfiles.Load(c.codeConfig.GenerationFolder)
	if err != nil {
		return err
	}
	c.ProjectState = projectState
	return nil
}

func (c *CodeGenerator) updateProjectState(update string) error {
	err := syncfiles.Update(c.codeConfig.GenerationFolder, update)
	if err != nil {
		return err
	}
	return c.loadProjectState()
}

func (c *CodeGenerator) promptModel() (string, error) {
	prompt := c.buildPrompt()
	c.Conversation.ResetMessages()
	responseMessage, err := c.Conversation.SendSystemMessage(prompt)
	if err != nil {
		return "", err
	}
	responseMessage.Serialize()
	responseContent := responseMessage.GetContent()
	return responseContent, nil
}

// Generates the code for the project based on the given directory state.
func (c *CodeGenerator) Generate() error {
	err := c.loadProjectState()
	if err != nil {
		return err
	}
	responseContent, err := c.promptModel()
	if err != nil {
		return err
	}
	return c.updateProjectState(responseContent)
}
