# brainstorm

## YAML-based 

### Processing

1. Start at the root of the repo and look for a `.mogent/config.yaml`
   file.
2. If found, read the config file and process any `include`
   directives.  An included file can be a local file path or a URL.
   Included files can themselves include other files.
3. If any definitions conflict, the last definition wins.
4. After processing all includes, the final definitions in the config
   file override any definitions in the included files.
5. After processing the config file, the tool will look for any `.md`
   files referenced in the compiled config and concatenate them into a
   single output file, in the order specified by the config.
6. Before writing the output file, the tool will process any Go templates in the
   .md files, substituting any variables defined in the config file or
   provided by the tool.
7. Write the final output to a file named `AGENTS.md` in the root of the repo.


### {repo}/.mogent/config.yaml:

All paths are relative to the same directory as the config file.

```
config:
    - include: http://github.com/ciwg/agents/main/ciwg.yaml
    - include: ~/.mogent.yaml

# following definitions are repo-specific and override any definitions in the included files
category:
    format:
        module:
            name: format
            source: repo-specific-format.md  # this path is {repo}/.mogent/repo-specific-format.md
    wire-lab:
        module:
            - name: sims
              source: sims.md  # this path is {repo}/.mogent/sims.md
            - name: pocs
              source: pocs.md  # this path is {repo}/.mogent/pocs.md
```

### templating

- Each included .md file is treated as a Go template.
- Dev can define variables in the config file. For example, if the
  config file defines a variable `project_name`, you can use 
  `{{ project_name }}` in your .md files to insert its value.
- Tool can define default variables, such as `repo_name`, `repo_url`,
  etc., that can be used in the templates. 
