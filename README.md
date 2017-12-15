# go-oai-link-extractor
A simple Go app that will, for a given OAI-PMH endpoint delivering oai_dc data, extract all the links in identifier tags and print them to stdout.

This currently only supports the `oai_dc` metadata prefix.

Using the PMH protocol to request data for specific time-frames and follow resumptions tokens this tool will retrieve all the links made available.

# TODO
Do we want to make this an all of time thing or have the supplied URL have the dates on it? But we should have an all of time approach for cases where all the links are needed. This needs to play into the idea of supporting repositories that will only deliver certain amounts of content at one time.

