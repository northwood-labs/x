package slices

import "testing"

var (
	// Not sorted, with duplicates.
	inputListStrings = []string{
		"apiaccountingconfig",
		"cloudcredentials",
		"activitylogalert",
		"containerinstance",
		"apprunnerservice",
		"bigdatainstance",
		"coldstorage",
		"batchpool",
		"apprunnerservice",
		"containerregistry",
		"accesslistflowlog",
		"bigdataserverlessnamespace",
		"cloudwatchdestination",
		"containernodegroup",
		"autoscalinggroup",
		"brokerinstance",
		"containerinstance",
		"cloudwatchdestination",
		"awsconfig",
		"azurepolicy",
		"accessanalyzer",
		"appstreamfleet",
		"containerinstance",
		"backupvault",
		"coderepository",
		"containerservice",
		"container",
		"apprunnerservice",
		"backupvault",
		"bigdataworkspace",
		"airflowenvironment",
		"cloudwatchdestination",
		"containerimage",
		"backupgateway",
		"backupvault",
		"artifactregistry",
		"apiaccountingconfig",
		"bigdataserverlessworkgroup",
		"contentdeliverynetwork",
		"automationaccount",
		"batchenvironment",
		"containerdeployment",
		"brokerinstance",
		"appserver",
		"apiaccountingconfig",
		"brokerinstance",
		"containerregistry",
		"autoscalinggroup",
		"bigdatasnapshot",
		"bigdataserverlessworkgroup",
		"backendservice",
		"botservice",
		"containerregistry",
		"autoscalinglaunchconfiguration",
		"containercluster",
		"bigdataserverlessworkgroup",
		"businessintelligencesubscription",
		"cassandratable",
		"buildproject",
	}

	// Sorted, without duplicates.
	outputListStrings = []string{
		"accessanalyzer",
		"accesslistflowlog",
		"activitylogalert",
		"airflowenvironment",
		"apiaccountingconfig",
		"apprunnerservice",
		"appserver",
		"appstreamfleet",
		"artifactregistry",
		"automationaccount",
		"autoscalinggroup",
		"autoscalinglaunchconfiguration",
		"awsconfig",
		"azurepolicy",
		"backendservice",
		"backupgateway",
		"backupvault",
		"batchenvironment",
		"batchpool",
		"bigdatainstance",
		"bigdataserverlessnamespace",
		"bigdataserverlessworkgroup",
		"bigdatasnapshot",
		"bigdataworkspace",
		"botservice",
		"brokerinstance",
		"buildproject",
		"businessintelligencesubscription",
		"cassandratable",
		"cloudcredentials",
		"cloudwatchdestination",
		"coderepository",
		"coldstorage",
		"container",
		"containercluster",
		"containerdeployment",
		"containerimage",
		"containerinstance",
		"containernodegroup",
		"containerregistry",
		"containerservice",
		"contentdeliverynetwork",
	}

	outputMapStrings = map[string]interface{}{
		"accessanalyzer":                   new(interface{}),
		"accesslistflowlog":                new(interface{}),
		"activitylogalert":                 new(interface{}),
		"airflowenvironment":               new(interface{}),
		"apiaccountingconfig":              new(interface{}),
		"apprunnerservice":                 new(interface{}),
		"appserver":                        new(interface{}),
		"appstreamfleet":                   new(interface{}),
		"artifactregistry":                 new(interface{}),
		"automationaccount":                new(interface{}),
		"autoscalinggroup":                 new(interface{}),
		"autoscalinglaunchconfiguration":   new(interface{}),
		"awsconfig":                        new(interface{}),
		"azurepolicy":                      new(interface{}),
		"backendservice":                   new(interface{}),
		"backupgateway":                    new(interface{}),
		"backupvault":                      new(interface{}),
		"batchenvironment":                 new(interface{}),
		"batchpool":                        new(interface{}),
		"bigdatainstance":                  new(interface{}),
		"bigdataserverlessnamespace":       new(interface{}),
		"bigdataserverlessworkgroup":       new(interface{}),
		"bigdatasnapshot":                  new(interface{}),
		"bigdataworkspace":                 new(interface{}),
		"botservice":                       new(interface{}),
		"brokerinstance":                   new(interface{}),
		"buildproject":                     new(interface{}),
		"businessintelligencesubscription": new(interface{}),
		"cassandratable":                   new(interface{}),
		"cloudcredentials":                 new(interface{}),
		"cloudwatchdestination":            new(interface{}),
		"coderepository":                   new(interface{}),
		"coldstorage":                      new(interface{}),
		"container":                        new(interface{}),
		"containercluster":                 new(interface{}),
		"containerdeployment":              new(interface{}),
		"containerimage":                   new(interface{}),
		"containerinstance":                new(interface{}),
		"containernodegroup":               new(interface{}),
		"containerregistry":                new(interface{}),
		"containerservice":                 new(interface{}),
		"contentdeliverynetwork":           new(interface{}),
	}

	inputListFloat32s = []float32{13.13, 1.1, 1.1, 5.5, 2.2, 8.8, 3.3}
	inputListFloat64s = []float64{13.13, 1.1, 1.1, 5.5, 2.2, 8.8, 3.3}
	inputListInt64s   = []int64{13, 1, 1, 5, 2, 8, 3}
	inputListInts     = []int{13, 1, 1, 5, 2, 8, 3}
	inputListUInt64s  = []uint64{13, 1, 1, 5, 2, 8, 3}
	inputListUInts    = []uint{13, 1, 1, 5, 2, 8, 3}

	outputListFloat32s = []float32{1.1, 2.2, 3.3, 5.5, 8.8, 13.13}
	outputListFloat64s = []float64{1.1, 2.2, 3.3, 5.5, 8.8, 13.13}
	outputListInt64s   = []int64{1, 2, 3, 5, 8, 13}
	outputListInts     = []int{1, 2, 3, 5, 8, 13}
	outputListUInt64s  = []uint64{1, 2, 3, 5, 8, 13}
	outputListUInts    = []uint{1, 2, 3, 5, 8, 13}
)

func TestSliceDedupeStrings(t *testing.T) {
	workingList := SliceDedupe(inputListStrings)
	outputList := outputListStrings

	if len(workingList) != len(outputList) {
		t.Error("input length is different from output length")
	}

	for i, v := range workingList {
		if v != outputList[i] {
			t.Errorf("input item `%v` does not match output item `%v`.", workingList[i], outputList[i])
		}
	}
}

func TestSliceDedupeFloat32s(t *testing.T) {
	workingList := SliceDedupe(inputListFloat32s)
	outputList := outputListFloat32s

	if len(workingList) != len(outputList) {
		t.Error("input length is different from output length")
	}

	for i, v := range workingList {
		if v != outputList[i] {
			t.Errorf("input item `%v` does not match output item `%v`.", workingList[i], outputList[i])
		}
	}
}

func TestSliceDedupeFloat64s(t *testing.T) {
	workingList := SliceDedupe(inputListFloat64s)
	outputList := outputListFloat64s

	if len(workingList) != len(outputList) {
		t.Error("input length is different from output length")
	}

	for i, v := range workingList {
		if v != outputList[i] {
			t.Errorf("input item `%v` does not match output item `%v`.", workingList[i], outputList[i])
		}
	}
}

func TestSliceDedupeInt64s(t *testing.T) {
	workingList := SliceDedupe(inputListInt64s)
	outputList := outputListInt64s

	if len(workingList) != len(outputList) {
		t.Error("input length is different from output length")
	}

	for i, v := range workingList {
		if v != outputList[i] {
			t.Errorf("input item `%v` does not match output item `%v`.", workingList[i], outputList[i])
		}
	}
}

func TestSliceDedupeInts(t *testing.T) {
	workingList := SliceDedupe(inputListInts)
	outputList := outputListInts

	if len(workingList) != len(outputList) {
		t.Error("input length is different from output length")
	}

	for i, v := range workingList {
		if v != outputList[i] {
			t.Errorf("input item `%v` does not match output item `%v`.", workingList[i], outputList[i])
		}
	}
}

func TestSliceDedupeUInt64s(t *testing.T) {
	workingList := SliceDedupe(inputListUInt64s)
	outputList := outputListUInt64s

	if len(workingList) != len(outputList) {
		t.Error("input length is different from output length")
	}

	for i, v := range workingList {
		if v != outputList[i] {
			t.Errorf("input item `%v` does not match output item `%v`.", workingList[i], outputList[i])
		}
	}
}

func TestSliceDedupeUInts(t *testing.T) {
	workingList := SliceDedupe(inputListUInts)
	outputList := outputListUInts

	if len(workingList) != len(outputList) {
		t.Error("input length is different from output length")
	}

	for i, v := range workingList {
		if v != outputList[i] {
			t.Errorf("input item `%v` does not match output item `%v`.", workingList[i], outputList[i])
		}
	}
}

func TestStringSliceToHashmap(t *testing.T) {
	workingMap := StringSliceToHashmap(outputListStrings)

	for _, v := range outputListStrings {
		if _, ok := workingMap[v]; !ok {
			t.Errorf("value `%v` is missing from resulting map.", v)
		}
	}
}
