//go:build darwin

package mpfirewall

func AddFirewallRule(ruleName string, exeFilePath string, port int) error {
	return nil
}

func AddFirewallRuleEx(ruleName string, exeFilePath string, port int, protocol int) error {
	return nil
}

func DeleteFirewallRule(ruleName string) error {
	return nil
}

func GetFirewallBlockedRules(exeFilePath string, port int) map[string]bool {
	return nil
}

func GetFirewallRuleLocalPort(ruleName string) (string, error) {
	return "", nil
}

func IsFirewallRuleExists(exeFilePath string, port int) (bool, error) {
	return false, nil
}
