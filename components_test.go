package discordgo

import (
	"encoding/json"
	"testing"
)

func TestModalChoiceComponentsMarshal(t *testing.T) {
	optional := false
	one := 1
	components := []MessageComponent{
		Label{Label: "Class", Component: RadioGroup{
			CustomID: "class",
			Options:  []RadioGroupOption{{Value: "a", Label: "A", Default: true}, {Value: "b", Label: "B"}},
			Required: &optional,
		}},
		Label{Label: "Days", Component: CheckboxGroup{
			CustomID:  "days",
			Options:   []CheckboxGroupOption{{Value: "mo", Label: "Mo", Description: "Monday"}},
			MinValues: &one,
			MaxValues: 1,
		}},
		Label{Label: "Agree", Component: Checkbox{CustomID: "agree", Default: true}},
	}

	got, err := json.Marshal(components)
	if err != nil {
		t.Fatal(err)
	}
	want := `[` +
		`{"label":"Class","component":{"custom_id":"class","options":[{"value":"a","label":"A","default":true},{"value":"b","label":"B"}],"required":false,"type":21},"type":18},` +
		`{"label":"Days","component":{"custom_id":"days","options":[{"value":"mo","label":"Mo","description":"Monday"}],"min_values":1,"max_values":1,"type":22},"type":18},` +
		`{"label":"Agree","component":{"custom_id":"agree","default":true,"type":23},"type":18}` +
		`]`
	if string(got) != want {
		t.Errorf("unexpected JSON\n got: %s\nwant: %s", got, want)
	}
}

func TestModalChoiceComponentsUnmarshal(t *testing.T) {
	raw := `{"custom_id":"modal","components":[
		{"id":1,"type":18,"component":{"id":2,"type":21,"custom_id":"class","value":"warrior"}},
		{"id":3,"type":18,"component":{"id":4,"type":21,"custom_id":"empty_radio","value":null}},
		{"id":5,"type":18,"component":{"id":6,"type":22,"custom_id":"days","values":["mo","di"]}},
		{"id":7,"type":18,"component":{"id":8,"type":23,"custom_id":"agree","value":true}},
		{"id":9,"type":18,"component":{"id":10,"type":23,"custom_id":"newsletter","value":false}}
	]}`

	var data ModalSubmitInteractionData
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		t.Fatal(err)
	}
	if len(data.Components) != 5 {
		t.Fatalf("expected 5 components, got %d", len(data.Components))
	}
	inner := func(i int) MessageComponent {
		label, ok := data.Components[i].(*Label)
		if !ok {
			t.Fatalf("component %d: expected *Label, got %T", i, data.Components[i])
		}
		return label.Component
	}

	radio, ok := inner(0).(*RadioGroup)
	if !ok || radio.CustomID != "class" || radio.Value == nil || *radio.Value != "warrior" {
		t.Errorf("radio group not parsed correctly: %#v", inner(0))
	}
	if empty, ok := inner(1).(*RadioGroup); !ok || empty.Value != nil {
		t.Errorf("radio group without selection should have nil value: %#v", inner(1))
	}
	group, ok := inner(2).(*CheckboxGroup)
	if !ok || len(group.Values) != 2 || group.Values[0] != "mo" || group.Values[1] != "di" {
		t.Errorf("checkbox group not parsed correctly: %#v", inner(2))
	}
	if cb, ok := inner(3).(*Checkbox); !ok || !cb.Value || cb.CustomID != "agree" {
		t.Errorf("checked checkbox not parsed correctly: %#v", inner(3))
	}
	if cb, ok := inner(4).(*Checkbox); !ok || cb.Value {
		t.Errorf("unchecked checkbox not parsed correctly: %#v", inner(4))
	}
}
