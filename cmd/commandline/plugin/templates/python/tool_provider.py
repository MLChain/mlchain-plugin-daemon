from typing import Any

from mlchain_plugin import ToolProvider
from mlchain_plugin.errors.tool import ToolProviderCredentialValidationError


class {{ .PluginName | SnakeToCamel }}Provider(ToolProvider):
    def _validate_credentials(self, credentials: dict[str, Any]) -> None:
        try:
            """
            IMPLEMENT YOUR VALIDATION HERE
            """
        except Exception as e:
            raise ToolProviderCredentialValidationError(str(e))
