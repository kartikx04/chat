
APPENV=integration-test go test -p 1 -v ./... 2>&1 | awk '
    /^=== RUN[[:space:]]/ {
        name = $0
        sub(/^=== RUN[[:space:]]+/, "", name)
        if (name ~ /\//) {
            printf "  \033[96m▶ RUN\033[0m  %s\n", name
        } else {
            printf "\033[1;96m▶ RUN\033[0m  %s\n", name
        }
        next
    }

    /^[[:space:]]*--- PASS:/ {
        name = $0
        sub(/^[[:space:]]*--- PASS:[[:space:]]*/, "", name)
        indent = ($0 ~ /^[[:space:]]{4,}--- PASS:/) ? "  " : ""
        printf "%s\033[1;95m✓ PASS\033[0m %s\n", indent, name
        next
    }

    /^[[:space:]]*--- FAIL:/ {
        name = $0
        sub(/^[[:space:]]*--- FAIL:[[:space:]]*/, "", name)
        printf "\033[1;31m✗ FAIL\033[0m %s\n", name
        next
    }

    /^[[:space:]]*--- SKIP:/ {
        name = $0
        sub(/^[[:space:]]*--- SKIP:[[:space:]]*/, "", name)
        printf "\033[33m↷ SKIP\033[0m %s\n", name
        next
    }

    /^PASS$/ {
        next
    }

    /^FAIL$/ {
        printf "\033[1;31m✗ TEST SUITE FAILED\033[0m\n"
        next
    }

    /^ok[[:space:]]/ {
        print "\033[90m" $0 "\033[0m"
        next
    }

    /^\?[[:space:]]/ {
        print "\033[90m" $0 "\033[0m"
        next
    }

    {
        print
    }
'
