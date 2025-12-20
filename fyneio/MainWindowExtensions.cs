using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Text;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Media;

namespace fyneio {
    public static class MainWindowExtensions {
        private static readonly Dictionary<string, string> ControlTypePrefixMap = new() {
            { "Button", "btn" },
            { "TextBox", "txt" },
            { "ComboBox", "cmb" },
            { "TextBlock", "lbl" },
            { "Label", "lbl" },
            { "CheckBox", "chk" },
            { "RadioButton", "rad" },
            { "Separator", "sep" },
            { "ListBox", "lst" },
            { "DataGrid", "grd" }
        };

        public static string ToConstantName(this string name) {
            var result = new StringBuilder();
            for (var i = 0; i < name.Length; i++) {
                var c = name[i];
                if (char.IsUpper(c) && i > 0) {
                    result.Append('_');
                }
                result.Append(char.ToUpper(c));
            }
            return result.ToString();
        }

        public static string ToEscapedString(this string input) {
            if (string.IsNullOrEmpty(input)) return "";
            return input
                .Replace("\\", "\\\\")
                .Replace("\"", "\\\"")
                .Replace("\n", "\\n")
                .Replace("\r", "\\r")
                .Replace("\t", "\\t");
        }

        public static string ToFileName(this string title) {
            var result = string.Empty;
            if (string.IsNullOrWhiteSpace(title)) {
                result = Guid.NewGuid().ToString("N")[..8];
            } else {
                var invalidChars = Path.GetInvalidFileNameChars();
                var safeName = new StringBuilder();
                foreach (var c in title.Trim()) {
                    if (!invalidChars.Contains(c)) {
                        safeName.Append(c);
                    }
                }
                result = safeName.ToString();
                if (string.IsNullOrWhiteSpace(result)) {
                    result = Guid.NewGuid().ToString("N")[..8];
                }
            }
            return result;
        }

        public static string ToGoFunctionName(this string structName) {
            return "New" + structName;
        }

        public static string ToGoStructName(this string windowTitle) {
            var result = string.Empty;
            if (string.IsNullOrWhiteSpace(windowTitle)) {
                result = Guid.NewGuid().ToString("N")[..8];
            } else {
                var cleaned = windowTitle.Trim()
                    .Replace(" ", "")
                    .Replace("_", "")
                    .Replace("-", "")
                    .Replace("(", "")
                    .Replace(")", "")
                    .Replace("（", "")
                    .Replace("）", "")
                    .Replace("·", "")
                    .Replace("/", "")
                    .Replace("\\", "")
                    .Replace(".", "")
                    .Replace(":", "")
                    .Replace("：", "");
                
                if (string.IsNullOrEmpty(cleaned)) {
                    result = Guid.NewGuid().ToString("N")[..8];
                } else {
                    var sb = new StringBuilder();
                    var nextUpper = true;
                    foreach (var c in cleaned) {
                        if (char.IsLetterOrDigit(c)) {
                            if (nextUpper) {
                                sb.Append(char.ToUpper(c));
                                nextUpper = false;
                            } else {
                                sb.Append(c);
                            }
                        } else {
                            nextUpper = true;
                        }
                    }
                    result = sb.ToString();
                }
            }
            return result;
        }

        public static bool IsAlphanumeric(this string text) {
            var result = true;
            foreach (var c in text) {
                if ((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == ' ' || c == '_' || c == '-') {
                    continue;
                }
                result = false;
                break;
            }
            return result;
        }

        public static string ToName(this string text, string prefix) {
            var cleaned = text.Trim()
                .Replace(":", "")
                .Replace("：", "")
                .Replace("(", "")
                .Replace(")", "")
                .Replace("（", "")
                .Replace("）", "")
                .Replace(" ", "")
                .Replace("\t", "")
                .Replace("-", "")
                .Replace("_", "")
                .Replace("·", "")
                .Replace("/", "")
                .Replace("\\", "")
                .Replace(".", "");

            string result;
            if (string.IsNullOrEmpty(cleaned)) {
                var guid2 = Guid.NewGuid().ToString("N")[..4];
                result = prefix + char.ToUpper(guid2[0]) + guid2[1..];
            } else {
                var sb = new StringBuilder();
                sb.Append(prefix);
                var nextUpper = true;
                foreach (var c in cleaned) {
                    if (char.IsLetterOrDigit(c)) {
                        if (nextUpper) {
                            sb.Append(char.ToUpper(c));
                            nextUpper = false;
                        } else {
                            sb.Append(c);
                        }
                    } else {
                        nextUpper = true;
                    }
                }
                result = sb.ToString();
            }
            return result;
        }

        public static string ToPrefix(this string controlType) {
            var normalizedType = controlType;
            if (normalizedType.StartsWith("CONTROL_TYPE_")) {
                normalizedType = normalizedType.Substring(13);
            }
            normalizedType = normalizedType.Replace("_", "");
            normalizedType = normalizedType switch {
                "BUTTON" => "Button",
                "LABEL" => "Label",
                "TEXTBOX" => "TextBox",
                "COMBOBOX" => "ComboBox",
                "CHECKBOX" => "CheckBox",
                "RADIOBUTTON" => "RadioButton",
                "SEPARATOR" => "Separator",
                "LISTBOX" => "ListBox",
                "DATAGRID" => "DataGrid",
                "TEXTBLOCK" => "TextBlock",
                _ => normalizedType
            };
            return ControlTypePrefixMap.TryGetValue(normalizedType, out var prefix) ? prefix : "ctrl";
        }

        public static string ToStandardXName(this string xName, string controlType) {
            var expectedPrefix = controlType.ToPrefix();
            string result;
            if (xName.StartsWith(expectedPrefix, StringComparison.OrdinalIgnoreCase)) {
                var rest = xName[expectedPrefix.Length..];
                if (rest.Length > 0) {
                    result = expectedPrefix + char.ToUpper(rest[0]) + rest[1..];
                } else {
                    var guid1 = Guid.NewGuid().ToString("N")[..6];
                    result = expectedPrefix + char.ToUpper(guid1[0]) + guid1[1..];
                }
            } else {
                var commonPrefixes = ControlTypePrefixMap.Values.ToArray();
                var found = false;
                string? tempResult = null;
                foreach (var prefix in commonPrefixes) {
                    if (xName.StartsWith(prefix, StringComparison.OrdinalIgnoreCase)) {
                        var rest = xName[prefix.Length..];
                        if (rest.Length > 0) {
                            tempResult = expectedPrefix + char.ToUpper(rest[0]) + rest[1..];
                        }
                        found = true;
                        break;
                    }
                }
                if (found && tempResult != null) {
                    result = tempResult;
                } else {
                    result = expectedPrefix + char.ToUpper(xName[0]) + xName[1..];
                }
            }
            return result;
        }

        public static string ToUniqueName(this string baseName, HashSet<string> usedNames) {
            string result;
            if (usedNames.Contains(baseName)) {
                var counter = 2;
                string newName;
                do {
                    newName = $"{baseName}{counter}";
                    counter++;
                } while (usedNames.Contains(newName));
                usedNames.Add(newName);
                result = newName;
            } else {
                usedNames.Add(baseName);
                result = baseName;
            }
            return result;
        }

        public static string ToVariableName(this string name) {
            if (string.IsNullOrEmpty(name)) return "ctrl";
            var firstChar = char.ToLower(name[0]);
            return firstChar + name[1..];
        }
    }
}