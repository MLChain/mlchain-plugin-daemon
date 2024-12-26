from collections.abc import Generator
from typing import Any

from mlchain_plugin import Tool
from mlchain_plugin.entities.tool import ToolInvokeMessage

class {{ .PluginName | SnakeToCamel }}Tool(Tool):
    def _invoke(self, tool_parameters: dict[str, Any]) -> Generator[ToolInvokeMessage]:
        yield self.create_json_message({
            "result": "Hello, world!"
        })
