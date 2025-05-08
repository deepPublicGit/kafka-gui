package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

type KafkaCluster struct {
	Name               string
	BootstrapServers   string
	TruststoreLocation string
	TruststorePassword string
	KeystoreLocation   string
	KeystorePassword   string
	SchemaRegistryURL  string
}

var clusters []KafkaCluster
var explorer *fyne.Container
var expandedClusters = make(map[string]bool)

var mockTopics = map[string][]string{
	"cluster1": {"topic1", "topic2"},
	"cluster2": {"topic3", "topic4"},
}

var mockTopicsData = map[string][]string{
	"topic1": {"data1", "data2"},
	"topic2": {"alpha", "beta", "gamma"},
	"topic3": {"data3", "data4"},
	"topic4": {"data5", "data6"},
}

func main() {
	a := app.New()
	w := a.NewWindow("Fyne GUI Layout Demo")

	// --- Right Side Panels ---
	rightTopPanel, rightBottomPanel := createRightPanels()
	explorer = createExplorer()

	// --- Menu Bar ---
	var updateExplorer func()
	addKafkaCluster := fyne.NewMenuItem("Add Kafka Cluster", func() {
		showKafkaClusterPopup(w, explorer, rightTopPanel, rightBottomPanel, updateExplorer)
	})
	fileMenu := fyne.NewMenu("File",
		addKafkaCluster,
		fyne.NewMenuItem("Quit", func() { a.Quit() }),
	)
	mainMenu := fyne.NewMainMenu(fileMenu)
	w.SetMainMenu(mainMenu)

	updateExplorer = func() {
		objs := []fyne.CanvasObject{widget.NewLabel("Kafka Clusters")}
		for _, c := range clusters {
			lbl := widget.NewLabel(c.Name)
			lbl.Wrapping = fyne.TextWrap(fyne.TextTruncateClip)
			isExpanded := expandedClusters[c.Name]
			clickable := &ClickableOverlay{OnTap: func(cluster KafkaCluster) func() {
				return func() {
					expandedClusters[cluster.Name] = !expandedClusters[cluster.Name]
					// Show cluster properties on cluster click
					rightTopPanel.Objects = []fyne.CanvasObject{
						widget.NewLabel("Bootstrap Servers: " + cluster.BootstrapServers),
						widget.NewLabel("Truststore Location: " + cluster.TruststoreLocation),
						widget.NewLabel("Truststore Password: " + cluster.TruststorePassword),
					}
					rightBottomPanel.Objects = []fyne.CanvasObject{
						widget.NewLabel("Keystore Location: " + cluster.KeystoreLocation),
						widget.NewLabel("Keystore Password: " + cluster.KeystorePassword),
						widget.NewLabel("Schema Registry URL: " + cluster.SchemaRegistryURL),
					}
					rightTopPanel.Refresh()
					rightBottomPanel.Refresh()
					updateExplorer()
				}
			}(c)}
			objs = append(objs, container.NewStack(lbl, clickable))

			if isExpanded {
				topics := mockTopics[c.Name]
				if topics == nil || len(topics) == 0 {
					topics = []string{"default topic"}
				}
				for _, topic := range topics {
					topicLbl := widget.NewLabel("    " + topic)
					topicClickable := &ClickableOverlay{OnTap: func(topic string) func() {
						return func() {
							// Show topic data in right panels
							vals := mockTopicsData[topic]
							if len(vals) == 0 {
								rightTopPanel.Objects = []fyne.CanvasObject{widget.NewLabel("No data for topic: " + topic)}
								rightBottomPanel.Objects = nil
							} else {
								rightTopPanel.Objects = []fyne.CanvasObject{widget.NewLabel("Topic: " + topic)}
								var bottom []fyne.CanvasObject
								for _, v := range vals {
									bottom = append(bottom, widget.NewLabel(v))
								}
								rightBottomPanel.Objects = bottom
							}
							rightTopPanel.Refresh()
							rightBottomPanel.Refresh()
						}
					}(topic)}
					objs = append(objs, container.NewStack(topicLbl, topicClickable))
				}
			}
		}
		explorer.Objects = objs
		explorer.Refresh()
	}
	// Initialize explorer
	updateExplorer()

	rightVSplit := container.NewVSplit(rightTopPanel, rightBottomPanel)
	rightVSplit.Offset = 0.5 // Split equally at start

	mainSplit := container.NewHSplit(explorer, rightVSplit)
	mainSplit.Offset = 0.2 // Explorer gets 20%

	w.SetContent(mainSplit)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}

func showKafkaClusterPopup(w fyne.Window, explorer *fyne.Container, rightTopPanel, rightBottomPanel *fyne.Container, updateExplorer func()) {
	clusterEntry := widget.NewEntry()
	bootstrapEntry := widget.NewEntry()
	truststoreLocEntry := widget.NewEntry()
	truststorePassEntry := widget.NewPasswordEntry()
	keystoreLocEntry := widget.NewEntry()
	keystorePassEntry := widget.NewPasswordEntry()
	schemaRegistryEntry := widget.NewEntry()
	var dlg *widget.PopUp // For modal dialog

	truststoreBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		filename, err := dialog.File().Title("Select Truststore").Load()
		if err == nil {
			truststoreLocEntry.SetText(filename)
		}
	})
	keystoreBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		filename, err := dialog.File().Title("Select Keystore").Load()
		if err == nil {
			keystoreLocEntry.SetText(filename)
		}
	})

	truststoreLocRow := container.NewBorder(nil, nil, nil, truststoreBtn, truststoreLocEntry)
	keystoreLocRow := container.NewBorder(nil, nil, nil, keystoreBtn, keystoreLocEntry)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Cluster Name", Widget: clusterEntry},
			{Text: "Bootstrap Servers", Widget: bootstrapEntry},
			{Text: "Truststore Location", Widget: truststoreLocRow},
			{Text: "Truststore Password", Widget: truststorePassEntry},
			{Text: "Keystore Location", Widget: keystoreLocRow},
			{Text: "Keystore Password", Widget: keystorePassEntry},
			{Text: "Schema Registry URL", Widget: schemaRegistryEntry},
		},
		SubmitText: "",
	}
	// Custom submit and cancel buttons
	submitBtn := widget.NewButton("Submit", func() {
		// Store the cluster
		cluster := KafkaCluster{
			Name:               clusterEntry.Text,
			BootstrapServers:   bootstrapEntry.Text,
			TruststoreLocation: truststoreLocEntry.Text,
			TruststorePassword: truststorePassEntry.Text,
			KeystoreLocation:   keystoreLocEntry.Text,
			KeystorePassword:   keystorePassEntry.Text,
			SchemaRegistryURL:  schemaRegistryEntry.Text,
		}
		clusters = append(clusters, cluster)
		// If cluster doesn't exist in mockTopics, add default
		if _, ok := mockTopics[cluster.Name]; !ok {
			mockTopics[cluster.Name] = []string{"default topic"}
		}
		updateExplorer()
		// Optionally select the newly added cluster
		rightTopPanel.Objects = []fyne.CanvasObject{
			widget.NewLabel("Bootstrap Servers: " + cluster.BootstrapServers),
			widget.NewLabel("Truststore Location: " + cluster.TruststoreLocation),
			widget.NewLabel("Truststore Password: " + cluster.TruststorePassword),
		}
		rightBottomPanel.Objects = []fyne.CanvasObject{
			widget.NewLabel("Keystore Location: " + cluster.KeystoreLocation),
			widget.NewLabel("Keystore Password: " + cluster.KeystorePassword),
			widget.NewLabel("Schema Registry URL: " + cluster.SchemaRegistryURL),
		}
		rightTopPanel.Refresh()
		rightBottomPanel.Refresh()
		updateExplorer()
		dlg.Hide()
	})
	cancelBtn := widget.NewButton("Cancel", func() { dlg.Hide() })
	buttonBar := container.NewHBox(
		layout.NewSpacer(),
		cancelBtn,
		submitBtn,
		layout.NewSpacer(),
	)

	popupWidth := w.Canvas().Size().Width * 0.4
	popupHeight := float32(420)
	margin := popupWidth * 0.05
	popup := container.NewVBox(
		widget.NewLabelWithStyle("Add Kafka Cluster", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(" "), // Top margin
		form,
		widget.NewLabel(" "),
		buttonBar,
		widget.NewLabel(" "), // Bottom margin
	)
	dlg = widget.NewModalPopUp(popup, w.Canvas())
	dlg.Resize(fyne.NewSize(popupWidth, popupHeight))
	dlg.Move(fyne.NewPos(
		(w.Canvas().Size().Width-popupWidth)/2+margin,
		(w.Canvas().Size().Height-popupHeight)/2-30, // vertical margin
	))
	dlg.Show()
}

func createExplorer() *fyne.Container {
	explorerItems := []fyne.CanvasObject{
		widget.NewLabel("Kafka Clusters"),
	}
	return container.NewVBox(explorerItems...)
}

func createRightPanels() (*fyne.Container, *fyne.Container) {
	rightTopPanel := container.NewVBox(
		widget.NewLabel("No Kafka Cluster Configured"),
	)
	rightBottomPanel := container.NewVBox()
	return rightTopPanel, rightBottomPanel
}

// ClickableOverlay implements a transparent clickable overlay for labels
// Used to make label rows clickable in the explorer
// Place this at the end of your file

type ClickableOverlay struct {
	widget.BaseWidget
	OnTap func()
}

func (c *ClickableOverlay) CreateRenderer() fyne.WidgetRenderer {
	rect := canvas.NewRectangle(color.Transparent)
	return widget.NewSimpleRenderer(rect)
}

func (c *ClickableOverlay) Tapped(_ *fyne.PointEvent) {
	if c.OnTap != nil {
		c.OnTap()
	}
}

func (c *ClickableOverlay) TappedSecondary(_ *fyne.PointEvent) {}
