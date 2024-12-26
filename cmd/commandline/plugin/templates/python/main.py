from mlchain_plugin import Plugin, MlchainPluginEnv

plugin = Plugin(MlchainPluginEnv(MAX_REQUEST_TIMEOUT=120))

if __name__ == '__main__':
    plugin.run()
