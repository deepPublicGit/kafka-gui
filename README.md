# Kafka Cluster GUI (Fyne)

This project is a Kafka Cluster management GUI built in Go using the [Fyne](https://fyne.io/) framework.

## Features & TODO

### ✅ Completed
- [x] **Basic Fyne GUI Layout** ([main.go](main.go))
- [x] **Add Kafka Cluster via Popup** ([main.go](main.go))
- [x] **Display Saved Clusters in Explorer** ([main.go](main.go))
- [x] **Expand Cluster to Show Topics (Mock Data)** ([main.go](main.go))
- [x] **Click Cluster to Show Properties** ([main.go](main.go))
- [x] **Click Topic to Show Mock Messages** ([main.go](main.go))
- [x] **Add/Remove Clusters in-memory** ([main.go](main.go))
- [x] **Resizable/Collapsible Explorer Panel** ([main.go](main.go))
- [x] **Password fields masked** ([main.go](main.go))

### ⏳ Pending / TODO
- [ ] **Test Button in Popup**
    - Validate cluster connection using [`kafkactl`](https://github.com/deviceinsight/kafkactl) (`get topics --config-file <temp>`)
    - Enable submit only if test passes
- [ ] **Write/Append Cluster Config to Local Config File**
- [ ] **On Cluster Click: Fetch Topics from kafkactl (not mock)**
    - Use context (cluster name) from config file
    - Cache topics per cluster
- [ ] **On Topic Click: Fetch Messages from kafkactl (not mock)**
    - Show top 10 messages in rightTopPanel
- [ ] **On Message Click: Show Full Message in rightBottomPanel**
- [ ] **Error Handling and User Feedback for all external commands**
- [ ] **Persist Clusters and Topics between sessions**
- [ ] **UI/UX improvements and refactoring**

## Usage
- Run `main.go` with Go 1.18+ and Fyne installed.
- To use real Kafka features, ensure [`kafkactl`](https://github.com/deviceinsight/kafkactl) is installed and available in your PATH.

## Links
- [Fyne Documentation](https://developer.fyne.io/)
- [kafkactl GitHub](https://github.com/deviceinsight/kafkactl)
