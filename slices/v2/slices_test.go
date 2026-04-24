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

	outputMapStrings = map[string]any{
		"accessanalyzer":                   new(any),
		"accesslistflowlog":                new(any),
		"activitylogalert":                 new(any),
		"airflowenvironment":               new(any),
		"apiaccountingconfig":              new(any),
		"apprunnerservice":                 new(any),
		"appserver":                        new(any),
		"appstreamfleet":                   new(any),
		"artifactregistry":                 new(any),
		"automationaccount":                new(any),
		"autoscalinggroup":                 new(any),
		"autoscalinglaunchconfiguration":   new(any),
		"awsconfig":                        new(any),
		"azurepolicy":                      new(any),
		"backendservice":                   new(any),
		"backupgateway":                    new(any),
		"backupvault":                      new(any),
		"batchenvironment":                 new(any),
		"batchpool":                        new(any),
		"bigdatainstance":                  new(any),
		"bigdataserverlessnamespace":       new(any),
		"bigdataserverlessworkgroup":       new(any),
		"bigdatasnapshot":                  new(any),
		"bigdataworkspace":                 new(any),
		"botservice":                       new(any),
		"brokerinstance":                   new(any),
		"buildproject":                     new(any),
		"businessintelligencesubscription": new(any),
		"cassandratable":                   new(any),
		"cloudcredentials":                 new(any),
		"cloudwatchdestination":            new(any),
		"coderepository":                   new(any),
		"coldstorage":                      new(any),
		"container":                        new(any),
		"containercluster":                 new(any),
		"containerdeployment":              new(any),
		"containerimage":                   new(any),
		"containerinstance":                new(any),
		"containernodegroup":               new(any),
		"containerregistry":                new(any),
		"containerservice":                 new(any),
		"contentdeliverynetwork":           new(any),
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

func TestDedupeStrings(t *testing.T) {
	workingList := Dedupe(inputListStrings)
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

func TestDedupeFloat32s(t *testing.T) {
	workingList := Dedupe(inputListFloat32s)
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

func TestDedupeFloat64s(t *testing.T) {
	workingList := Dedupe(inputListFloat64s)
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

func TestDedupeInt64s(t *testing.T) {
	workingList := Dedupe(inputListInt64s)
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

func TestDedupeInts(t *testing.T) {
	workingList := Dedupe(inputListInts)
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

func TestDedupeUInt64s(t *testing.T) {
	workingList := Dedupe(inputListUInt64s)
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

func TestDedupeUInts(t *testing.T) {
	workingList := Dedupe(inputListUInts)
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

type filterSubstrTestUser struct {
	User string
	Key  string
}

func TestFilterSubstrStructSlice(t *testing.T) {
	input := []filterSubstrTestUser{
		{User: "alice", Key: "team-red"},
		{User: "bob", Key: "team-blue"},
		{User: "charlie", Key: "ops"},
	}

	workingList := FilterSubstr(input, "blue", func(u filterSubstrTestUser) string {
		return u.User + " " + u.Key
	})

	if len(workingList) != 1 {
		t.Fatalf("unexpected filtered length: got %d want %d", len(workingList), 1)
	}

	if workingList[0].User != "bob" {
		t.Errorf("unexpected user in filtered output: got %q want %q", workingList[0].User, "bob")
	}
}

func TestFilterSubstrStringSlice(t *testing.T) {
	input := []string{"alpha", "bravo", "charlie"}

	workingList := FilterSubstr(input, "ha", func(s string) string {
		return s
	})

	if len(workingList) != 2 {
		t.Fatalf("unexpected filtered length: got %d want %d", len(workingList), 2)
	}

	if workingList[0] != "alpha" {
		t.Errorf("unexpected first match: got %q want %q", workingList[0], "alpha")
	}

	if workingList[1] != "charlie" {
		t.Errorf("unexpected second match: got %q want %q", workingList[1], "charlie")
	}
}

func TestFilterRegexStructSlice(t *testing.T) {
	input := []filterSubstrTestUser{
		{User: "alice", Key: "team-red"},
		{User: "bob", Key: "team-blue"},
		{User: "charlie", Key: "ops"},
	}

	workingList := FilterRegex(input, "^bob|blue$", func(u filterSubstrTestUser) string {
		return u.User + " " + u.Key
	})

	if len(workingList) != 1 {
		t.Fatalf("unexpected filtered length: got %d want %d", len(workingList), 1)
	}

	if workingList[0].User != "bob" {
		t.Errorf("unexpected user in filtered output: got %q want %q", workingList[0].User, "bob")
	}
}

func TestFilterRegexStringSlice(t *testing.T) {
	input := []string{"alpha", "bravo", "charlie", "delta"}

	workingList := FilterRegex(input, "^(char|del)", func(s string) string {
		return s
	})

	if len(workingList) != 2 {
		t.Fatalf("unexpected filtered length: got %d want %d", len(workingList), 2)
	}

	if workingList[0] != "charlie" {
		t.Errorf("unexpected first match: got %q want %q", workingList[0], "charlie")
	}

	if workingList[1] != "delta" {
		t.Errorf("unexpected second match: got %q want %q", workingList[1], "delta")
	}
}
