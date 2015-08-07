
install: nr_agent_sdk

AGENT_SDK_VERSION=v0.16.2.0-beta.x86_64

nr_agent_sdk-${AGENT_SDK_VERSION}.tar.gz:
	wget http://download.newrelic.com/agent_sdk/nr_agent_sdk-${AGENT_SDK_VERSION}.tar.gz

nr_agent_sdk: nr_agent_sdk-${AGENT_SDK_VERSION}.tar.gz
	tar xzf $<
	ln -s nr_agent_sdk-${AGENT_SDK_VERSION} nr_agent_sdk
