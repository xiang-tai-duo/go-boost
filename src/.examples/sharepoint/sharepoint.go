// Package main
// File:        https://github.com/xiang-tai-duo/go-boost/blob/master/.examples/sharepoint/sharepoint.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: Sharepoint usage example
// --------------------------------------------------------------------------------
package main

import (
	"fmt"

	__sharepoint "github.com/xiang-tai-duo/go-boost/sharepoint"
)

func main() {
	// Create SharePoint client with client ID and client secret
	siteURL := "https://your-sharepoint-site.sharepoint.com/sites/yoursite"
	clientID := "your-client-id"
	clientSec := "your-client-secret"
	err := error(nil)
	sp := __sharepoint.New(siteURL, clientID, clientSec)

	// Get site URL
	var site string
	if site, err = sp.GetSiteURL(); err == nil {
		fmt.Printf("Site URL: %s\n", site)
	} else {
		fmt.Printf("Error getting site URL: %v\n", err)
	}

	// Get authentication token
	var token string
	if token, err = sp.GetToken(); err == nil {
		fmt.Printf("Token: %s\n", token)
	} else {
		fmt.Printf("Error getting token: %v\n", err)
	}

	// Example 3: File operations
	// Upload file
	localFilePath := "local-file.txt"
	remoteFilePath := "/Shared Documents/local-file.txt"
	if err = sp.UploadFile(localFilePath, remoteFilePath); err == nil {
		fmt.Println("File uploaded successfully")
	} else {
		fmt.Printf("Error uploading file: %v\n", err)
	}

	// Download file
	downloadPath := "downloaded-file.txt"
	if err = sp.DownloadFile(remoteFilePath, downloadPath); err == nil {
		fmt.Println("File downloaded successfully")
	} else {
		fmt.Printf("Error downloading file: %v\n", err)
	}

	// List files
	folderPath := "/Shared Documents"
	var files []string
	if files, err = sp.ListFiles(folderPath); err == nil {
		fmt.Printf("Files in %s: %v\n", folderPath, files)
	} else {
		fmt.Printf("Error listing files: %v\n", err)
	}

	// Create folder
	newFolderPath := "/Shared Documents/New Folder"
	if err = sp.CreateFolder(newFolderPath); err == nil {
		fmt.Println("Folder created successfully")
	} else {
		fmt.Printf("Error creating folder: %v\n", err)
	}

	// Delete file
	if err = sp.DeleteFile(remoteFilePath); err == nil {
		fmt.Println("File deleted successfully")
	} else {
		fmt.Printf("Error deleting file: %v\n", err)
	}

	// Example 4: List operations
	listName := "Documents"

	// Get list items
	var items []map[string]interface{}
	if items, err = sp.GetListItems(listName); err == nil {
		fmt.Printf("List items: %v\n", items)
	} else {
		fmt.Printf("Error getting list items: %v\n", err)
	}

	// Add list item
	newItem := map[string]interface{}{
		"Title":       "New Document",
		"Description": "This is a new document",
	}
	if err = sp.AddListItem(listName, newItem); err == nil {
		fmt.Println("List item added successfully")
	} else {
		fmt.Printf("Error adding list item: %v\n", err)
	}

	// Update list item
	itemID := 1
	updatedItem := map[string]interface{}{
		"Title":       "Updated Document",
		"Description": "This document has been updated",
	}
	if err = sp.UpdateListItem(listName, itemID, updatedItem); err == nil {
		fmt.Println("List item updated successfully")
	} else {
		fmt.Printf("Error updating list item: %v\n", err)
	}

	// Delete list item
	if err = sp.DeleteListItem(listName, itemID); err == nil {
		fmt.Println("List item deleted successfully")
	} else {
		fmt.Printf("Error deleting list item: %v\n", err)
	}
}
