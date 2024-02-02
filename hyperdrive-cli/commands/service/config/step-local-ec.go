package config

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/nodeset-org/hyperdrive/shared/types"
)

func createLocalEcStep(wiz *wizard, currentStep int, totalSteps int) *choiceWizardStep {
	// Make lists of clients that good and bad
	goodClients := []*types.ParameterOption[types.ExecutionClient]{}
	badClients := []*types.ParameterOption[types.ExecutionClient]{}
	for _, client := range wiz.md.Config.Hyperdrive.LocalExecutionConfig.ExecutionClient.Options {
		if !strings.HasPrefix(client.Name, "*") {
			goodClients = append(goodClients, client)
		} else {
			badClients = append(badClients, client)
		}
	}

	randomDesc := strings.Builder{}
	randomDesc.WriteString("Select a client randomly to help promote the diversity of the Ethereum Chain. We recommend you do this unless you have a strong reason to pick a specific client.")
	if len(badClients) > 0 {
		randomDesc.WriteString("\n\n[orange]NOTE: The following clients are currently overrepresented on the Ethereum network (\"supermajority\" clients):\n\t")
		labels := []string{}
		for _, client := range badClients {
			labels = append(labels, strings.TrimPrefix(client.Name, "*"))
		}
		randomDesc.WriteString(strings.Join(labels, ", "))
		randomDesc.WriteString("\nWe recommend choosing different clients for the health of the network. Please see https://clientdiversity.org/ to learn more.")
	}

	// Create the button names and descriptions from the config
	clients := wiz.md.Config.Hyperdrive.LocalExecutionConfig.ExecutionClient.Options
	clientNames := []string{"Random (Recommended)"}
	clientDescriptions := []string{randomDesc.String()}
	for _, client := range clients {
		clientNames = append(clientNames, client.Name)
		clientDescriptions = append(clientDescriptions, client.Description)
	}

	helperText := "Please select the Execution Client you would like to use.\n\nHighlight each one to see a brief description of it, or go to https://clientdiversity.org/ to learn more about them."

	show := func(modal *choiceModalLayout) {
		wiz.md.setPage(modal.page)
		modal.focus(0) // Catch-all for safety

		if !wiz.md.isNew {
			var ecName string
			for _, option := range wiz.md.Config.Hyperdrive.LocalExecutionConfig.ExecutionClient.Options {
				if option.Value == wiz.md.Config.Hyperdrive.LocalExecutionConfig.ExecutionClient.Value {
					ecName = option.Name
					break
				}
			}
			for i, clientName := range clientNames {
				if ecName == clientName {
					modal.focus(i)
					break
				}
			}
		}
	}

	done := func(buttonIndex int, buttonLabel string) {
		if buttonIndex == 0 {
			wiz.md.pages.RemovePage(randomBnPrysmID)
			wiz.md.pages.RemovePage(randomBnID)
			selectRandomEC(goodClients, wiz, currentStep, totalSteps)
		} else {
			buttonLabel = strings.TrimSpace(buttonLabel)
			selectedClient := types.ExecutionClient_Unknown
			for _, client := range wiz.md.Config.Hyperdrive.LocalExecutionConfig.ExecutionClient.Options {
				if client.Name == buttonLabel {
					selectedClient = client.Value
					break
				}
			}
			if selectedClient == types.ExecutionClient_Unknown {
				panic(fmt.Sprintf("Local EC selection buttons didn't match any known clients, buttonLabel = %s\n", buttonLabel))
			}
			wiz.md.Config.Hyperdrive.LocalExecutionConfig.ExecutionClient.Value = selectedClient
			wiz.bnLocalModal.show()
		}
	}

	back := func() {
		wiz.modeModal.show()
	}

	return newChoiceStep(
		wiz,
		currentStep,
		totalSteps,
		helperText,
		clientNames,
		clientDescriptions,
		100,
		"Execution Client > Selection",
		DirectionalModalVertical,
		show,
		done,
		back,
		"step-ec-local",
	)
}

// Get a random execution client
func selectRandomEC(goodOptions []*types.ParameterOption[types.ExecutionClient], wiz *wizard, currentStep int, totalSteps int) {
	// Get system specs
	//totalMemoryGB := memory.TotalMemory() / 1024 / 1024 / 1024
	//isLowPower := (totalMemoryGB < 15 || runtime.GOARCH == "arm64")

	// Filter out the clients based on system specs
	filteredClients := []types.ExecutionClient{}
	for _, clientOption := range goodOptions {
		client := clientOption.Value
		switch client {
		default:
			filteredClients = append(filteredClients, client)
		}
	}

	// Select a random client
	rand.Seed(time.Now().UnixNano())
	selectedClient := filteredClients[rand.Intn(len(filteredClients))]
	wiz.md.Config.Hyperdrive.LocalExecutionConfig.ExecutionClient.Value = selectedClient

	// Show the selection page
	wiz.executionLocalRandomModal = createRandomEcStep(wiz, currentStep, totalSteps, goodOptions)
	wiz.executionLocalRandomModal.show()
}
