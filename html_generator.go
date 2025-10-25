package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func generateHtml(planData interface{}) string {
	// Parse the plan data to extract resource changes and drift
	planMap, ok := planData.(map[string]interface{})
	if !ok {
		return generateErrorHtml("Invalid plan data format")
	}

	// Extract resource changes
	resourceChanges := extractResourceChanges(planMap)
	driftCount := countDriftChanges(planMap)

	// Generate HTML content
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Terraform Plan</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        h1 {
            color: #2c3e50;
            border-bottom: 2px solid #3498db;
            padding-bottom: 10px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .source-link {
            font-size: 14px;
            font-weight: normal;
            color: #3498db;
            text-decoration: none;
        }
        .source-link:hover {
            text-decoration: underline;
        }
        .promo-message {
            text-align: center;
            margin: 20px 0;
            font-style: italic;
        }
        .promo-link {
            color: #3498db;
            text-decoration: none;
            font-weight: bold;
            font-style: normal;
        }
        .promo-link:hover {
            text-decoration: underline;
        }
        .section-header-row {
            display: flex;
            justify-content: space-between;
            align-items: center;
            width: 100%;
        }
        .section-description {
            font-size: 14px;
            font-style: italic;
            color: #6c757d;
            margin-bottom: 15px;
        }
        .section {
            margin: 20px 0;
            padding: 15px;
            background-color: #ecf0f1;
            border-radius: 5px;
        }
        .resource-item {
            margin: 10px 0;
            padding: 10px;
            background-color: white;
            border-radius: 3px;
            border-left: 4px solid #3498db;
        }
        .create { border-left-color: #27ae60; }
        .update { border-left-color: #f39c12; }
        .delete { border-left-color: #e74c3c; }
        .no-op { border-left-color: #95a5a6; }
        .drift { border-left-color: #e67e22; }
        .replace { border-left-color: #95a5a6; }
        .action {
            font-weight: bold;
            padding: 2px 8px;
            border-radius: 3px;
            color: white;
            font-size: 12px;
        }
        .action-create { background-color: #27ae60; }
        .action-update { background-color: #f39c12; }
        .action-delete { background-color: #e74c3c; }
        .action-no-op { background-color: #95a5a6; }
        .action-drift { background-color: #e67e22; }
        .action-replace { background: linear-gradient(90deg, #e74c3c 40%, #27ae60 60%); }
        .resource-address {
            font-family: monospace;
            font-weight: bold;
            color: #2c3e50;
        }
        .resource-type {
            color: #7f8c8d;
            font-size: 14px;
        }
        .resource-attributes {
            margin-top: 10px;
            padding: 10px;
            background-color: #f8f9fa;
            border-radius: 3px;
            font-family: monospace;
            font-size: 12px;
        }
        .diff-container {
            display: flex;
            gap: 10px;
            margin-top: 10px;
        }
        .diff-container pre {
            white-space: pre-wrap;
            word-wrap: break-word;
            margin: 0;
            padding: 5px;
            background-color: rgba(0,0,0,0.05);
            border-radius: 3px;
        }
        .collapsible {
            cursor: pointer;
            user-select: none;
            display: flex;
            align-items: center;
            gap: 8px;
        }
        .collapsible:hover {
            background-color: #f0f0f0;
        }
        .collapsible::before {
            content: "▼";
            font-size: 12px;
            transition: transform 0.2s;
            flex-shrink: 0;
        }
        .collapsible.collapsed::before {
            content: "▶";
        }
        .collapsible-content {
            overflow: hidden;
            transition: opacity 0.3s ease-out, max-height 0.3s ease-out;
        }
        .collapsible-content.collapsed {
            max-height: 0;
            opacity: 0;
        }
        .collapsible-content:not(.collapsed) {
            max-height: none;
            opacity: 1;
        }
        .diff-column {
            flex: 1;
            padding: 10px;
            border-radius: 3px;
        }
        .diff-before {
            background-color: #f8d7da;
            border-left: 3px solid #dc3545;
        }
        .diff-after {
            background-color: #d4edda;
            border-left: 3px solid #28a745;
        }
        .diff-header {
            font-weight: bold;
            margin-bottom: 10px;
            color: #495057;
        }
        .attribute-item {
            margin: 5px 0;
            padding: 3px 0;
            border-bottom: 1px solid #e9ecef;
        }
        .attribute-key {
            font-weight: bold;
            color: #495057;
        }
        .attribute-value {
            color: #6c757d;
            margin-left: 10px;
        }
        .attribute-changed {
            background-color: #fff3cd;
            border-left: 3px solid #ffc107;
            padding-left: 8px;
        }
        .attribute-added {
            background-color: #d4edda;
            border-left: 3px solid #28a745;
            padding-left: 8px;
        }
        .attribute-removed {
            background-color: #f8d7da;
            border-left: 3px solid #dc3545;
            padding-left: 8px;
        }
        .attribute-computed {
            background-color: #e2e3e5;
            border-left: 3px solid #6c757d;
            padding-left: 8px;
        }
        .summary {
            display: flex;
            gap: 20px;
            margin-bottom: 20px;
        }
        .summary-item {
            flex: 1;
            text-align: center;
            padding: 15px;
            background-color: white;
            border-radius: 5px;
        }
        .summary-number {
            font-size: 24px;
            font-weight: bold;
            color: #2c3e50;
        }
        .summary-label {
            color: #7f8c8d;
            font-size: 14px;
        }
    </style>
    <script>
        function toggleCollapsible(element) {
            const content = element.nextElementSibling;
            element.classList.toggle('collapsed');
            content.classList.toggle('collapsed');
        }
        
        // Make individual resource items collapsed by default, but keep main sections open
        document.addEventListener('DOMContentLoaded', function() {
            const collapsibles = document.querySelectorAll('.collapsible');
            collapsibles.forEach(function(element) {
                // Check if this is a main section (Resource Changes or Resource Drift)
                const isMainSection = element.querySelector('h2') !== null;
                
                if (!isMainSection) {
                    // Only collapse individual resource items, not main sections
                    element.classList.add('collapsed');
                    const content = element.nextElementSibling;
                    if (content) {
                        content.classList.add('collapsed');
                    }
                } else {
                    // Check if main section has no resource items
                    const section = element.closest('.section');
                    const resourceItems = section.querySelectorAll('.resource-item');
                    if (resourceItems.length === 0) {
                        element.classList.add('collapsed');
                        const content = element.nextElementSibling;
                        if (content) {
                            content.classList.add('collapsed');
                        }
                    }
                }
            });
        });
    </script>
</head>
<body>
    <div class="container">
        <h1>Terraform Plan</h1>
                
        <div class="section">
            <div class="collapsible" onclick="toggleCollapsible(this)">
                <div class="section-header-row">
                    <h2>Resource Changes (` + fmt.Sprintf("%d", len(resourceChanges)) + ` total)</h2>
                    <p class="section-description">Terraform will apply these changes to your resources</p>
                </div>
            </div>
            <div class="collapsible-content">
                
                ` + generateResourceChangesHtml(resourceChanges) + `
            </div>
        </div>
        
        <div class="section">
            <div class="collapsible collapsed" onclick="toggleCollapsible(this)">
            <div class="section-header-row">
                    <h2>Resource Drift (` + fmt.Sprintf("%d", driftCount) + ` total)</h2>
                    <p class="section-description">Terraform will update these resources to match your configuration</p>
                </div>
            </div>
            <div class="collapsible-content collapsed">
                <div style="text-align: center; padding: 40px 20px; background-color: #f8f9fa; border-radius: 8px; margin: 20px 0;">
                    <h3 style="color: #2c3e50; margin-bottom: 15px;">Want to see detailed drift information?</h3>
                    <p style="color: #6c757d; margin-bottom: 20px; font-size: 16px;">
                        Get comprehensive drift analysis, cost estimates, and advanced Terraform insights with CloudVIC.
                    </p>
                    <a href="https://cloudvic.com" style="display: inline-block; background-color: #3498db; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; font-weight: bold; transition: background-color 0.3s;">
                        Subscribe to CloudVIC
                    </a>
                </div>
            </div>
        </div>
    </div>
    <div class="promo-message">
        Want to visualize your Terraform plan and state changes over time, enriched with pricing API data and drift detection?<br>
        <a href="https://cloudvic.com" class="promo-link">Try CloudVIC</a>
    </div>
</body>
</html>`

	return html
}

func extractResourceChanges(planMap map[string]interface{}) []map[string]interface{} {
	var changes []map[string]interface{}
	var driftAddresses []string

	// First, collect addresses from drift that have delete actions
	if driftChanges, ok := planMap["resource_drift"].([]interface{}); ok {
		for _, change := range driftChanges {
			if changeMap, ok := change.(map[string]interface{}); ok {
				actions := getActions(changeMap)
				if len(actions) > 0 && actions[0] == "delete" {
					address := getString(changeMap, "address")
					driftAddresses = append(driftAddresses, address)
				}
			}
		}
	}

	// Then process resource changes, marking as replace if they also appear in drift
	if resourceChanges, ok := planMap["resource_changes"].([]interface{}); ok {
		for _, change := range resourceChanges {
			if changeMap, ok := change.(map[string]interface{}); ok {
				// Filter out no-op changes
				actions := getActions(changeMap)
				if len(actions) > 0 && actions[0] != "no-op" {
					address := getString(changeMap, "address")

					// Check if this resource also appears in drift (indicating replace)
					if contains(driftAddresses, address) && actions[0] == "create" {
						// Mark as replace operation
						changeMap["_is_replace"] = true
					}
					changes = append(changes, changeMap)
				}
			}
		}
	}

	return changes
}

func countDriftChanges(planMap map[string]interface{}) int {
	var replaceAddresses []string

	// First, collect addresses that are being replaced (appear in both changes and drift)
	if driftChanges, ok := planMap["resource_drift"].([]interface{}); ok {
		for _, change := range driftChanges {
			if changeMap, ok := change.(map[string]interface{}); ok {
				actions := getActions(changeMap)
				if len(actions) > 0 && actions[0] == "delete" {
					address := getString(changeMap, "address")

					// Check if this address also appears in resource_changes as create
					if resourceChanges, ok := planMap["resource_changes"].([]interface{}); ok {
						for _, rc := range resourceChanges {
							if rcMap, ok := rc.(map[string]interface{}); ok {
								rcAddress := getString(rcMap, "address")
								rcActions := getActions(rcMap)
								if rcAddress == address && len(rcActions) > 0 && rcActions[0] == "create" {
									replaceAddresses = append(replaceAddresses, address)
									break
								}
							}
						}
					}
				}
			}
		}
	}

	// Count drift changes excluding replace operations
	driftCount := 0
	if driftChanges, ok := planMap["resource_drift"].([]interface{}); ok {
		for _, change := range driftChanges {
			if changeMap, ok := change.(map[string]interface{}); ok {
				address := getString(changeMap, "address")
				if !contains(replaceAddresses, address) {
					driftCount++
				}
			}
		}
	}

	return driftCount
}

func generateResourceChangesHtml(changes []map[string]interface{}) string {
	if len(changes) == 0 {
		return "<p>No resource changes detected.</p>"
	}

	var html strings.Builder
	html.WriteString("<div>")

	for _, change := range changes {
		address := getString(change, "address")
		actions := getActions(change)

		// Check if this is a replace operation
		var displayActions []string
		if _, isReplace := change["_is_replace"]; isReplace {
			displayActions = []string{"replace"}
		} else {
			displayActions = actions
		}

		// Get change details
		changeDetails := getChangeDetails(change)

		html.WriteString(fmt.Sprintf(`
			<div class="resource-item %s">
				<div class="collapsible" onclick="toggleCollapsible(this)">
					<div>%s</div>
					<div class="resource-address">%s</div>
				</div>
				<div class="collapsible-content">
					<div class="resource-attributes">
						%s
					</div>
				</div>
			</div>`,
			getActionClass(displayActions[0]),
			formatActions(displayActions),
			address,
			changeDetails))
	}

	html.WriteString("</div>")
	return html.String()
}

func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func getActions(change map[string]interface{}) []string {
	if changeData, ok := change["change"].(map[string]interface{}); ok {
		if actions, ok := changeData["actions"].([]interface{}); ok {
			var result []string
			for _, action := range actions {
				if actionStr, ok := action.(string); ok {
					result = append(result, actionStr)
				}
			}
			return result
		}
	}
	return []string{"no-op"}
}

func getActionClass(action string) string {
	switch action {
	case "create":
		return "create"
	case "update":
		return "update"
	case "delete":
		return "delete"
	case "replace":
		return "replace"
	default:
		return "no-op"
	}
}

func formatActions(actions []string) string {
	var result []string
	for _, action := range actions {
		class := "action-" + getActionClass(action)
		result = append(result, fmt.Sprintf(`<span class="action %s">%s</span>`, class, strings.ToUpper(action)))
	}
	return strings.Join(result, " ")
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func getChangeDetails(change map[string]interface{}) string {
	var details strings.Builder

	if changeData, ok := change["change"].(map[string]interface{}); ok {
		actions := getActions(change)
		action := actions[0]

		// Check if this is a replace operation
		if _, isReplace := change["_is_replace"]; isReplace {
			// For replace operations, show both before and after
			details.WriteString("<div class='attribute-item'><span class='attribute-key'>Resource Replacement:</span></div>")

			// Show what's being removed
			if before, ok := changeData["before"].(map[string]interface{}); ok {
				details.WriteString("<div class='attribute-item'><span class='attribute-key'>Current State (to be removed):</span></div>")
				details.WriteString(formatAttributes(before, "attribute-removed"))
			}

			// Show what's being created
			if after, ok := changeData["after"].(map[string]interface{}); ok {
				details.WriteString("<div class='attribute-item'><span class='attribute-key'>New State (to be created):</span></div>")
				details.WriteString(formatAttributes(after, "attribute-added"))
			}

			// Show computed/unknown fields
			if afterUnknown, ok := changeData["after_unknown"].(map[string]interface{}); ok {
				details.WriteString("<div class='attribute-item'><span class='attribute-key'>Computed Fields:</span></div>")
				details.WriteString(formatUnknownAttributes(afterUnknown, "attribute-computed"))
			}
		} else if action == "create" {
			// For creates, show the "after" values and "after_unknown" fields
			if after, ok := changeData["after"].(map[string]interface{}); ok {
				details.WriteString("<div class='attribute-item'><span class='attribute-key'>New Resource:</span></div>")
				details.WriteString(formatAttributes(after, "attribute-added"))
			}

			// Also show computed/unknown fields
			if afterUnknown, ok := changeData["after_unknown"].(map[string]interface{}); ok {
				details.WriteString("<div class='attribute-item'><span class='attribute-key'>Computed Fields:</span></div>")
				details.WriteString(formatUnknownAttributes(afterUnknown, "attribute-computed"))
			}
		} else if action == "delete" {
			// For deletes, show only the "before" values
			if before, ok := changeData["before"].(map[string]interface{}); ok {
				details.WriteString("<div class='attribute-item'><span class='attribute-key'>Resource to Delete:</span></div>")
				details.WriteString(formatAttributes(before, "attribute-removed"))
			}
		} else if action == "update" {
			// For updates, show side-by-side diff of only changed fields
			before, beforeOk := changeData["before"].(map[string]interface{})
			after, afterOk := changeData["after"].(map[string]interface{})

			if beforeOk && afterOk {
				changedFields := getChangedFields(before, after)
				if len(changedFields) > 0 {
					details.WriteString("<div class='diff-container'>")
					details.WriteString("<div class='diff-column'>")
					details.WriteString("<div class='diff-header'>Before</div>")
					details.WriteString(formatChangedFields(changedFields, before, "attribute-removed"))
					details.WriteString("</div>")

					details.WriteString("<div class='diff-column'>")
					details.WriteString("<div class='diff-header'>After</div>")
					details.WriteString(formatChangedFields(changedFields, after, "attribute-added"))
					details.WriteString("</div>")
					details.WriteString("</div>")
				}
			}
		}
	}

	return details.String()
}

func formatAttributes(attrs map[string]interface{}, cssClass string) string {
	var result strings.Builder

	for key, value := range attrs {
		valueStr := formatValue(value)
		if valueStr != "" {
			// Check if the value contains newlines (JSON formatting)
			if strings.Contains(valueStr, "\n") {
				result.WriteString(fmt.Sprintf(`
					<div class="attribute-item %s">
						<span class="attribute-key">%s:</span>
						<pre class="attribute-value">%s</pre>
					</div>`, cssClass, key, valueStr))
			} else {
				result.WriteString(fmt.Sprintf(`
					<div class="attribute-item %s">
						<span class="attribute-key">%s:</span>
						<span class="attribute-value">%s</span>
					</div>`, cssClass, key, valueStr))
			}
		}
	}

	return result.String()
}

func formatUnknownAttributes(attrs map[string]interface{}, cssClass string) string {
	var result strings.Builder

	for key, value := range attrs {
		// For unknown attributes, we show "known after creation" as the value
		// The value from after_unknown is typically a boolean indicating if it's unknown
		if unknownBool, ok := value.(bool); ok && unknownBool {
			result.WriteString(fmt.Sprintf(`
				<div class="attribute-item %s">
					<span class="attribute-key">%s:</span>
					<span class="attribute-value">known after creation</span>
				</div>`, cssClass, key))
		} else if _, ok := value.(map[string]interface{}); ok {
			// Handle nested unknown objects
			result.WriteString(fmt.Sprintf(`
				<div class="attribute-item %s">
					<span class="attribute-key">%s:</span>
					<span class="attribute-value">known after creation</span>
				</div>`, cssClass, key))
		} else if _, ok := value.([]interface{}); ok {
			// Handle unknown arrays
			result.WriteString(fmt.Sprintf(`
				<div class="attribute-item %s">
					<span class="attribute-key">%s:</span>
					<span class="attribute-value">known after creation</span>
				</div>`, cssClass, key))
		}
	}

	return result.String()
}

func getChangedFields(before, after map[string]interface{}) []string {
	var changedFields []string

	// Check all keys in both before and after
	allKeys := make(map[string]bool)
	for key := range before {
		allKeys[key] = true
	}
	for key := range after {
		allKeys[key] = true
	}

	for key := range allKeys {
		beforeVal, beforeExists := before[key]
		afterVal, afterExists := after[key]

		// Field was added
		if !beforeExists && afterExists {
			changedFields = append(changedFields, key)
		} else if beforeExists && !afterExists {
			// Field was removed
			changedFields = append(changedFields, key)
		} else if beforeExists && afterExists && !valuesEqual(beforeVal, afterVal) {
			// Field was changed
			changedFields = append(changedFields, key)
		}
	}

	return changedFields
}

func valuesEqual(a, b interface{}) bool {
	// Handle nil cases
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Type assertion and comparison
	switch va := a.(type) {
	case string:
		if vb, ok := b.(string); ok {
			return va == vb
		}
	case float64:
		if vb, ok := b.(float64); ok {
			return va == vb
		}
	case bool:
		if vb, ok := b.(bool); ok {
			return va == vb
		}
	case []interface{}:
		if vb, ok := b.([]interface{}); ok {
			if len(va) != len(vb) {
				return false
			}
			for i, v := range va {
				if !valuesEqual(v, vb[i]) {
					return false
				}
			}
			return true
		}
	case map[string]interface{}:
		if vb, ok := b.(map[string]interface{}); ok {
			if len(va) != len(vb) {
				return false
			}
			for k, v := range va {
				if !valuesEqual(v, vb[k]) {
					return false
				}
			}
			return true
		}
	}

	return false
}

func formatChangedFields(changedFields []string, data map[string]interface{}, cssClass string) string {
	var result strings.Builder

	for _, key := range changedFields {
		if value, exists := data[key]; exists {
			valueStr := formatValue(value)
			// Show the value even if it's empty (empty strings are valid in diffs)

			// Check if the value contains newlines (JSON formatting)
			if strings.Contains(valueStr, "\n") {
				result.WriteString(fmt.Sprintf(`
					<div class="attribute-item %s">
						<span class="attribute-key">%s:</span>
						<pre class="attribute-value">%s</pre>
					</div>`, cssClass, key, valueStr))
			} else {
				result.WriteString(fmt.Sprintf(`
					<div class="attribute-item %s">
						<span class="attribute-key">%s:</span>
						<span class="attribute-value">%s</span>
					</div>`, cssClass, key, valueStr))
			}
		}
	}

	return result.String()
}

func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		// Check if the string contains JSON
		if isJSONString(v) {
			// Try to parse and reformat the JSON
			var jsonData interface{}
			if err := json.Unmarshal([]byte(v), &jsonData); err == nil {
				jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
				if err == nil {
					return string(jsonBytes)
				}
			}
		}
		if len(v) > 100 {
			return v[:100] + "..."
		}
		return v
	case float64:
		return fmt.Sprintf("%.0f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case []interface{}:
		if len(v) == 0 {
			return "[]"
		}
		// Format as pretty-printed JSON
		jsonBytes, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return fmt.Sprintf("[%d items]", len(v))
		}
		return string(jsonBytes)
	case map[string]interface{}:
		// Format as pretty-printed JSON
		jsonBytes, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return fmt.Sprintf("{%d fields}", len(v))
		}
		return string(jsonBytes)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// isJSONString checks if a string contains JSON
func isJSONString(s string) bool {
	if len(s) == 0 {
		return false
	}
	// Check if it starts and ends with JSON-like characters
	trimmed := strings.TrimSpace(s)
	return (strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}")) ||
		(strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]"))
}

func generateErrorHtml(message string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><title>Error</title></head>
<body><h1>Error</h1><p>%s</p></body>
</html>`, message)
}
