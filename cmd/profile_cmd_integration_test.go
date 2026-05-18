package cmd

import (
	"os"
	"testing"

	"envcrypt/crypto"
)

func TestRunProfileAdd_Overwrite(t *testing.T) {
	dir := setupProfileTestDir(t)
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	if err := runProfileAdd([]string{"prod", ".env.prod"}); err != nil {
		t.Fatalf("first add: %v", err)
	}
	if err := runProfileAdd([]string{"prod", ".env.prod.v2"}); err != nil {
		t.Fatalf("second add: %v", err)
	}
	store, err := crypto.LoadProfileStore(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	p, ok := store.GetProfile("prod")
	if !ok {
		t.Fatal("expected profile to exist")
	}
	if p.EnvFile != ".env.prod.v2" {
		t.Errorf("expected updated env file, got %s", p.EnvFile)
	}
}

func TestRunProfileAdd_MultipleProfiles(t *testing.T) {
	dir := setupProfileTestDir(t)
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	profiles := [][]string{
		{"dev", ".env.dev"},
		{"staging", ".env.staging"},
		{"prod", ".env.prod"},
	}
	for _, p := range profiles {
		if err := runProfileAdd(p); err != nil {
			t.Fatalf("add %s: %v", p[0], err)
		}
	}
	store, err := crypto.LoadProfileStore(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(store.Profiles) != 3 {
		t.Errorf("expected 3 profiles, got %d", len(store.Profiles))
	}
	if err := runProfileList([]string{}); err != nil {
		t.Fatalf("list: %v", err)
	}
}

func TestRunProfileRemove_ThenList(t *testing.T) {
	dir := setupProfileTestDir(t)
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	runProfileAdd([]string{"dev", ".env.dev"})
	runProfileAdd([]string{"prod", ".env.prod"})
	runProfileRemove([]string{"dev"})

	store, err := crypto.LoadProfileStore(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if _, ok := store.GetProfile("dev"); ok {
		t.Error("expected 'dev' to be removed")
	}
	if _, ok := store.GetProfile("prod"); !ok {
		t.Error("expected 'prod' to remain")
	}
	if err := runProfileList([]string{}); err != nil {
		t.Fatalf("list: %v", err)
	}
}
