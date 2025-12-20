#pragma warning disable IDE0003

using System.Diagnostics;
using System.IO;
using System.Reflection;
using System.Text;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Media;

namespace fyneio {
    public partial class MainWindow : Window {
        private const string CONTROL_TYPE_BUTTON = "BUTTON";
        private const string CONTROL_TYPE_CHECK_BOX = "CHECK_BOX";
        private const string CONTROL_TYPE_COMBO_BOX = "COMBO_BOX";
        private const string CONTROL_TYPE_DATA_GRID = "DATA_GRID";
        private const string CONTROL_TYPE_LABEL = "LABEL";
        private const string CONTROL_TYPE_LIST_BOX = "LIST_BOX";
        private const string CONTROL_TYPE_RADIO_BUTTON = "RADIO_BUTTON";
        private const string CONTROL_TYPE_SEPARATOR = "SEPARATOR";
        private const string CONTROL_TYPE_TEXT_BLOCK = "TEXT_BLOCK";
        private const string CONTROL_TYPE_TEXT_BOX = "TEXT_BOX";

        private static readonly Dictionary<string, string> CONTROL_TYPE_MAP = new() {
            { "Button", CONTROL_TYPE_BUTTON },
            { "CheckBox", CONTROL_TYPE_CHECK_BOX },
            { "ComboBox", CONTROL_TYPE_COMBO_BOX },
            { "DataGrid", CONTROL_TYPE_DATA_GRID },
            { "Label", CONTROL_TYPE_LABEL },
            { "ListBox", CONTROL_TYPE_LIST_BOX },
            { "RadioButton", CONTROL_TYPE_RADIO_BUTTON },
            { "Separator", CONTROL_TYPE_SEPARATOR },
            { "TextBlock", CONTROL_TYPE_TEXT_BLOCK },
            { "TextBox", CONTROL_TYPE_TEXT_BOX }
        };

        private readonly string[] TARGET_CONTROL_TYPES = [
            "Button",
            "TextBox",
            "TextBlock",
            "Label",
            "ComboBox",
            "CheckBox",
            "RadioButton",
            "ListBox",
            "DataGrid",
            "Separator"
        ];
        private readonly List<ControlInfo> CONTROLS = [];
        private readonly Dictionary<ControlInfo, string> CONTROL_NAMES = [];

        private class ControlInfo {
            public string Content { get; set; } = string.Empty;
            public string Name { get; set; } = string.Empty;
            public string Placeholder { get; set; } = string.Empty;
            public string? SelectedItem { get; set; }
            public string Type { get; set; } = string.Empty;
            public double Height { get; set; }
            public double Left { get; set; }
            public double Top { get; set; }
            public double Width { get; set; }
            public bool IsMultiLine { get; set; }
            public bool IsReadOnly { get; set; }
            public List<string>? Items { get; set; }
        }

        public MainWindow() {
            InitializeComponent();
            this.WindowStartupLocation = WindowStartupLocation.CenterScreen;
        }

        private void Window_Loaded(object sender, RoutedEventArgs e) {
            WalkControls(this);
            CreateControlNames();
            CreateFyneCode();
        }

        private string CreateBaseName(ControlInfo ctrl) {
            string result;
            if (!string.IsNullOrEmpty(ctrl.Name)) {
                result = ctrl.Name.ToStandardXName(ctrl.Type);
            } else {
                var text = ctrl.Content;
                if (!string.IsNullOrWhiteSpace(text) && text.IsAlphanumeric()) {
                    result = text.ToName(ctrl.Type.ToPrefix());
                } else {
                    var guid = Guid.NewGuid().ToString("N")[..6];
                    result = ctrl.Type.ToPrefix() + char.ToUpper(guid[0]) + guid[1..];
                }
            }
            return result;
        }

        private void CreateControlNames() {
            var usedNames = new HashSet<string>();
            foreach (var ctrl in CONTROLS) {
                var baseName = CreateBaseName(ctrl);
                var uniqueName = baseName.ToUniqueName(usedNames);
                CONTROL_NAMES[ctrl] = uniqueName;
            }
        }

        private void CreateFyneCode() {
            var structName = this.Title.ToGoStructName();
            var functionName = structName.ToGoFunctionName();
            var sb = new StringBuilder();
            sb.AppendLine($"package main");
            sb.AppendLine($"import (");
            sb.AppendLine($"    _ \"embed\"");
            sb.AppendLine($"    _v2 \"fyne.io/fyne/v2\"");
            sb.AppendLine($"    _app \"fyne.io/fyne/v2/app\"");
            sb.AppendLine($"    _container \"fyne.io/fyne/v2/container\"");
            sb.AppendLine($"    _theme \"fyne.io/fyne/v2/theme\"");
            sb.AppendLine($"    _widget \"fyne.io/fyne/v2/widget\"");
            sb.AppendLine($"    \"github.com/xiang-tai-duo/go-boost/fyneio\"");
            sb.AppendLine($")");
            sb.AppendLine($"");
            sb.AppendLine($"//goland:noinspection GoUnusedConst,GoSnakeCaseUsage,SpellCheckingInspection");
            sb.AppendLine($"const (");
            sb.AppendLine($"    SOURCEHANSANS_FONT_FILE_NAME = \"SourceHanSans-VF.ttf\"");
            sb.AppendLine($")");
            sb.AppendLine($"");
            sb.AppendLine($"//goland:noinspection GoSnakeCaseUsage,GoUnusedGlobalVariable,SpellCheckingInspection");
            sb.AppendLine($"var (");
            sb.AppendLine($"    //go:embed resources/fonts/SourceHanSans-VF.ttf");
            sb.AppendLine($"    SOURCEHANSANS_VF []byte");
            sb.AppendLine($")");
            sb.AppendLine($"");
            sb.AppendLine($"//goland:noinspection SpellCheckingInspection");
            sb.AppendLine($"type {structName} struct {{");
            sb.AppendLine($"    app    _v2.App");
            sb.AppendLine($"    window _v2.Window");
            sb.AppendLine($"}}");
            sb.AppendLine($"");
            sb.AppendLine($"//goland:noinspection SpellCheckingInspection");
            sb.AppendLine($"func {functionName}() *{structName} {{");
            sb.AppendLine($"    app := _app.New()");
            sb.AppendLine($"    theme := fyneio.NewTheme(_theme.VariantDark)");
            sb.AppendLine($"    theme.ThemeFont = fyneio.LoadFont(SOURCEHANSANS_FONT_FILE_NAME, SOURCEHANSANS_VF)");
            sb.AppendLine($"    theme.SizeNameText = 12");
            sb.AppendLine($"    theme.SizeNameHeadingText = 24");
            sb.AppendLine($"    theme.SizeNameSubHeadingText = 20");
            sb.AppendLine($"    theme.SizeNameCaptionText = 14");
            sb.AppendLine($"    app.Settings().SetTheme(theme)");
            sb.AppendLine($"    return &{structName}{{");
            sb.AppendLine($"        app: app,");
            sb.AppendLine($"    }}");
            sb.AppendLine($"}}");
            sb.AppendLine($"");
            sb.AppendLine($"func (instance *{structName}) Show() {{");
            sb.AppendLine($"    instance.window = instance.app.NewWindow(\"{this.Title}\")");
            sb.AppendLine($"    instance.window.Resize(_v2.NewSize({Math.Round(this.ActualWidth)}, {Math.Round(this.ActualHeight)}))");
            sb.AppendLine($"    instance.window.CenterOnScreen()");
            sb.AppendLine($"    instance.window.SetFixedSize(true)");
            sb.AppendLine($"    if win, ok := interface{{}}(instance.window).(interface{{ SetResizable(bool) }}); ok {{");
            sb.AppendLine($"        win.SetResizable(false)");
            sb.AppendLine($"    }}");
            sb.AppendLine($"    window := _container.NewWithoutLayout()");
            foreach (var ctrl in CONTROLS) {
                var varName = CONTROL_NAMES[ctrl].ToVariableName();
                var left = Math.Round(ctrl.Left);
                var top = Math.Round(ctrl.Top);
                var width = Math.Round(ctrl.Width);
                var height = Math.Round(ctrl.Height);
                switch (ctrl.Type) {
                    case CONTROL_TYPE_BUTTON:
                        sb.AppendLine($"{varName} := _widget.NewButton(\"{ctrl.Content.ToEscapedString()}\", nil)");
                        sb.AppendLine($"{varName}.Resize(_v2.NewSize({width}, {height}))");
                        sb.AppendLine($"{varName}.Move(_v2.NewPos({left}, {top}))");
                        sb.AppendLine($"window.Add({varName})");
                        break;
                    case CONTROL_TYPE_TEXT_BLOCK:
                    case CONTROL_TYPE_LABEL:
                        sb.AppendLine($"{varName} := fyneio.NewLabel(\"{ctrl.Content.ToEscapedString()}\", _theme.VariantDark)");
                        sb.AppendLine($"{varName}.Resize(_v2.NewSize({width}, {height}))");
                        sb.AppendLine($"{varName}.Move(_v2.NewPos({left}, {top}))");
                        sb.AppendLine($"window.Add({varName})");
                        break;
                    case CONTROL_TYPE_TEXT_BOX:
                        if (ctrl.IsMultiLine) {
                            sb.AppendLine($"{varName} := _widget.NewMultiLineEntry()");
                        } else {
                            sb.AppendLine($"{varName} := _widget.NewEntry()");
                        }
                        sb.AppendLine($"{varName}.Resize(_v2.NewSize({width}, {height}))");
                        sb.AppendLine($"{varName}.Move(_v2.NewPos({left}, {top}))");
                        if (!string.IsNullOrEmpty(ctrl.Content)) {
                            sb.AppendLine($"{varName}.SetText(\"{ctrl.Content.ToEscapedString()}\")");
                        }
                        if (!string.IsNullOrEmpty(ctrl.Placeholder)) {
                            sb.AppendLine($"{varName}.SetPlaceHolder(\"{ctrl.Placeholder.ToEscapedString()}\")");
                        }
                        sb.AppendLine($"window.Add({varName})");
                        break;
                    case CONTROL_TYPE_COMBO_BOX:
                        var items = ctrl.Items != null && ctrl.Items.Count > 0
                            ? string.Join(", ", ctrl.Items.Select(x => $"\"{x.ToEscapedString()}\""))
                            : string.Empty;
                        sb.AppendLine($"{varName} := _widget.NewSelect([]string{{{items}}}, nil)");
                        sb.AppendLine($"{varName}.Resize(_v2.NewSize({width}, {height}))");
                        sb.AppendLine($"{varName}.Move(_v2.NewPos({left}, {top}))");
                        if (!string.IsNullOrEmpty(ctrl.SelectedItem)) {
                            sb.AppendLine($"{varName}.SetSelected(\"{ctrl.SelectedItem.ToEscapedString()}\")");
                        }
                        sb.AppendLine($"window.Add({varName})");
                        break;
                    case CONTROL_TYPE_CHECK_BOX:
                        sb.AppendLine($"{varName} := _widget.NewCheck(\"{ctrl.Content.ToEscapedString()}\", nil)");
                        sb.AppendLine($"{varName}.Resize(_v2.NewSize({width}, {height}))");
                        sb.AppendLine($"{varName}.Move(_v2.NewPos({left}, {top}))");
                        sb.AppendLine($"window.Add({varName})");
                        break;
                    case CONTROL_TYPE_RADIO_BUTTON:
                        sb.AppendLine($"{varName} := _widget.NewRadioGroup([]string{{\"{ctrl.Content.ToEscapedString()}\"}}, nil)");
                        sb.AppendLine($"{varName}.Resize(_v2.NewSize({width}, {height}))");
                        sb.AppendLine($"{varName}.Move(_v2.NewPos({left}, {top}))");
                        sb.AppendLine($"window.Add({varName})");
                        break;
                    case CONTROL_TYPE_SEPARATOR:
                        sb.AppendLine($"{varName} := _widget.NewSeparator()");
                        sb.AppendLine($"{varName}.Resize(_v2.NewSize({width}, {height}))");
                        sb.AppendLine($"{varName}.Move(_v2.NewPos({left}, {top}))");
                        sb.AppendLine($"window.Add({varName})");
                        break;
                    default:
                        sb.AppendLine($"// Unhandled control type: {ctrl.Type}");
                        break;
                }
            }
            sb.AppendLine($"instance.window.SetContent(window)");
            sb.AppendLine($"instance.window.ShowAndRun()");
            sb.AppendLine($"}}");
            var exePath = Assembly.GetExecutingAssembly().Location;
            var exeDirectory = Path.GetDirectoryName(exePath)!;
            var fileName = this.Title.ToFileName();
            var outputPath = Path.Combine(exeDirectory, $"{fileName}.go");
            File.WriteAllText(outputPath, sb.ToString(), Encoding.UTF8);
            var binFontsPath = Path.Combine(exeDirectory, "resources", "fonts");
            if (!Directory.Exists(binFontsPath)) {
                Directory.CreateDirectory(binFontsPath);
            }
            var assembly = Assembly.GetExecutingAssembly();
            var fontResourceName = "fyneio.res.fonts.SourceHanSans-VF.ttf";
            var licenseResourceName = "fyneio.res.fonts.LICENSE.txt";
            var targetFontPath = Path.Combine(binFontsPath, "SourceHanSans-VF.ttf");
            using (var stream = assembly.GetManifestResourceStream(fontResourceName)) {
                if (stream != null) {
                    using var fileStream = new FileStream(targetFontPath, FileMode.Create, FileAccess.Write);
                    stream.CopyTo(fileStream);
                }
            }
            var targetLicensePath = Path.Combine(binFontsPath, "LICENSE.txt");
            using (var stream = assembly.GetManifestResourceStream(licenseResourceName)) {
                if (stream != null) {
                    using var fileStream = new FileStream(targetLicensePath, FileMode.Create, FileAccess.Write);
                    stream.CopyTo(fileStream);
                }
            }
            Application.Current.Shutdown();
        }

        private void WalkControls(DependencyObject parent) {
            for (var i = 0; i < VisualTreeHelper.GetChildrenCount(parent); i++) {
                var child = VisualTreeHelper.GetChild(parent, i);
                if (child is FrameworkElement element) {
                    var controlType = element.GetType().Name;
                    if (Array.Exists(TARGET_CONTROL_TYPES, t => t == controlType)) {
                        var relativePosition = element.TransformToAncestor(this).Transform(new Point(0, 0));
                        var left = relativePosition.X;
                        var top = relativePosition.Y;
                        var width = element.ActualWidth;
                        var height = element.ActualHeight;
                        var controlInfo = new ControlInfo {
                            Name = element.Name ?? string.Empty,
                            Type = CONTROL_TYPE_MAP.TryGetValue(controlType, out var mappedType) ? mappedType : controlType,
                            Left = left,
                            Top = top,
                            Width = width,
                            Height = height
                        };
                        switch (element) {
                            case Button btn:
                                controlInfo.Content = btn.Content?.ToString() ?? string.Empty;
                                break;
                            case TextBlock tb:
                                controlInfo.Content = tb.Text ?? string.Empty;
                                break;
                            case Label lbl:
                                controlInfo.Content = lbl.Content?.ToString() ?? string.Empty;
                                break;
                            case TextBox txt:
                                controlInfo.Content = txt.Text ?? string.Empty;
                                controlInfo.IsReadOnly = txt.IsReadOnly;
                                controlInfo.IsMultiLine = txt.TextWrapping != TextWrapping.NoWrap;
                                controlInfo.Placeholder = txt.Tag?.ToString() ?? string.Empty;
                                break;
                            case ComboBox cmb:
                                controlInfo.Items = cmb.Items.Cast<object>().Select(x => x.ToString() ?? string.Empty).ToList();
                                controlInfo.SelectedItem = cmb.SelectedItem?.ToString();
                                break;
                            case CheckBox chk:
                                controlInfo.Content = chk.Content?.ToString() ?? string.Empty;
                                break;
                            case RadioButton rad:
                                controlInfo.Content = rad.Content?.ToString() ?? string.Empty;
                                break;
                        }
                        CONTROLS.Add(controlInfo);
                    }
                    if (element is Button or CheckBox or RadioButton) {
                        continue;
                    }
                    WalkControls(element);
                }
            }
        }
    }
}