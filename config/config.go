package config

import (
	"fmt"
	"os"
	"regexp"
	"slices"

	"gopkg.in/yaml.v3"
)

type Regexp regexp.Regexp

func (r *Regexp) UnmarshalText(b []byte) error {
	re, err := regexp.Compile(string(b))
	if err != nil {
		return err
	}
	*r = Regexp(*re)

	return nil
}

func (r *Regexp) MarshalText() ([]byte, error) {
	if r != nil {
		return []byte((*regexp.Regexp)(r).String()), nil
	}
	return nil, nil
}

var allowedSkipQuestions = []string{
	"scope",
	"body",
	"breaking",
	"footer",
}

type Config struct {
	Types                []TypeValue   `yaml:"types"`
	Messages             Messages      `yaml:"messages,omitempty"`
	SkipQuestions        SkipQuestions `yaml:"skip_questions,omitempty"`
	AllowBreakingChanges []string      `yaml:"allow_breaking_changes,omitempty"`
	TicketNumber         TicketNumber  `yaml:"ticket_number,omitempty"`
}

type TypeValue struct {
	Value string `yaml:"value"`
	Name  string `yaml:"name"`
}

func (t TypeValue) String() string {
	return t.Name
}

type Messages struct {
	Type            string `yaml:"type,omitempty"`
	Scope           string `yaml:"scope,omitempty"`
	TicketNumber    string `yaml:"ticket_number,omitempty"`
	Subject         string `yaml:"subject,omitempty"`
	Body            string `yaml:"body,omitempty"`
	BreakingConfirm string `yaml:"breaking_confirm,omitempty"`
	BreakingMessage string `yaml:"breaking_message,omitempty"`
	Footer          string `yaml:"footer,omitempty"`
	ConfirmCommit   string `yaml:"confirm_commit,omitempty"`
}

type SkipQuestions []string

type TicketNumber struct {
	Enable         bool           `yaml:"enable,omitempty"`
	Required       bool           `yaml:"required,omitempty"`
	Prefix         string         `yaml:"prefix,omitempty"`
	MatchPattern   *Regexp        `yaml:"match_pattern,omitempty"`
	FromBranchName FromBranchName `yaml:"from_branch_name,omitempty"`
}

type FromBranchName struct {
	Enable        bool    `yaml:"enable,omitempty"`
	ExtractRegexp *Regexp `yaml:"extract_regexp,omitempty"`
}

func LoadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		return nil, err
	}

	for _, s := range cfg.SkipQuestions {
		if !slices.Contains(allowedSkipQuestions, s) {
			return nil, fmt.Errorf("invalid skip question: %s", s)
		}
	}
	return &cfg, nil
}
