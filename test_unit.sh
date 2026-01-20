#!/bin/bash
# Go Unit Test Runner with pretty output

echo "🧪 Go Unit Tests"
echo "============================"

PASS=0
FAIL=0

# Run tests and capture output
OUTPUT=$(go test ./... -v 2>&1)

# Parse and format output
while IFS= read -r line; do
  # Test suite header
  if [[ "$line" =~ ^===\ RUN\ +(.+)$ ]]; then
    TEST_NAME="${BASH_REMATCH[1]}"
    # Only show top-level tests (not subtests)
    if [[ ! "$TEST_NAME" =~ "/" ]]; then
      echo ""
      echo "📋 $TEST_NAME"
    fi
  # Passed test
  elif [[ "$line" =~ ^---\ PASS:\ +(.+) ]]; then
    TEST_NAME=$(echo "${BASH_REMATCH[1]}" | cut -d' ' -f1)
    # Only show top-level tests
    if [[ ! "$TEST_NAME" =~ "/" ]]; then
      echo "   ✅ $TEST_NAME passed"
      ((PASS++))
    fi
  # Failed test
  elif [[ "$line" =~ ^---\ FAIL:\ +(.+) ]]; then
    TEST_NAME=$(echo "${BASH_REMATCH[1]}" | cut -d' ' -f1)
    if [[ ! "$TEST_NAME" =~ "/" ]]; then
      echo "   ❌ $TEST_NAME failed"
      ((FAIL++))
    fi
  # Package result
  elif [[ "$line" =~ ^ok[[:space:]]+(.+)[[:space:]]+([0-9.]+s)$ ]]; then
    PKG="${BASH_REMATCH[1]}"
    TIME="${BASH_REMATCH[2]}"
    # Extract just the package name
    PKG_SHORT=$(echo "$PKG" | sed 's/wine-shop-api\///')
    echo "   ✅ Package: $PKG_SHORT ($TIME)"
  elif [[ "$line" =~ ^FAIL[[:space:]]+(.+) ]]; then
    PKG="${BASH_REMATCH[1]}"
    PKG_SHORT=$(echo "$PKG" | sed 's/wine-shop-api\///')
    echo "   ❌ Package failed: $PKG_SHORT"
  fi
done <<< "$OUTPUT"

echo ""
echo "============================"
echo "Results: $PASS passed, $FAIL failed"
echo "============================"

if [ $FAIL -gt 0 ]; then
  exit 1
fi
