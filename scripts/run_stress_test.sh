#!/bin/bash
# Long-running stress test for memory leak detection
# 
# Usage:
#   ./run_stress_test.sh [duration]
#   ./run_stress_test.sh 5h    # Run for 5 hours (default)
#   ./run_stress_test.sh 30m   # Run for 30 minutes
#
# The test results are saved to:
#   - stress_test_output.log    (test output)
#   - stress_memory_leak.log    (memory metrics CSV)
#
# To monitor progress in real-time:
#   tail -f stress_test_output.log
#
# After the test, analyze results:
#   go run cmd/analyze_memory_log.go stress_memory_leak.log --plot

set -e

# Configuration
DURATION=${1:-5h}
OUTPUT_DIR="stress_test_results"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

echo "══════════════════════════════════════════════════════════════"
echo "Gig Memory Leak Stress Test"
echo "══════════════════════════════════════════════════════════════"
echo "Duration: $DURATION"
echo "Output directory: $OUTPUT_DIR"
echo "Timestamp: $TIMESTAMP"
echo "══════════════════════════════════════════════════════════════"

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Calculate timeout (duration + 10 minutes)
case $DURATION in
    *h)
        HOURS=${DURATION%h}
        TIMEOUT_MINUTES=$((HOURS * 60 + 10))
        ;;
    *m)
        MINUTES=${DURATION%m}
        TIMEOUT_MINUTES=$((MINUTES + 10))
        ;;
    *)
        echo "Invalid duration format. Use format like 5h or 30m"
        exit 1
        ;;
esac

echo "Timeout: ${TIMEOUT_MINUTES}m"
echo ""

# Change to project root
cd "$(dirname "$0")/.."

# Clean up old log files
rm -f stress_memory_leak.log

# Run the test
echo "Starting test at $(date)"
echo "To monitor progress: tail -f $OUTPUT_DIR/stress_test_$TIMESTAMP.log"
echo ""

STRESS_DURATION=$DURATION go test ./tests/ \
    -run TestStress_MemoryLeak \
    -v \
    -count=1 \
    -timeout "${TIMEOUT_MINUTES}m" \
    2>&1 | tee "$OUTPUT_DIR/stress_test_$TIMESTAMP.log"

# Save memory log
if [ -f stress_memory_leak.log ]; then
    mv stress_memory_leak.log "$OUTPUT_DIR/stress_memory_leak_$TIMESTAMP.log"
    echo ""
    echo "══════════════════════════════════════════════════════════════"
    echo "Test completed at $(date)"
    echo "══════════════════════════════════════════════════════════════"
    echo "Output files:"
    echo "  - $OUTPUT_DIR/stress_test_$TIMESTAMP.log"
    echo "  - $OUTPUT_DIR/stress_memory_leak_$TIMESTAMP.log"
    echo ""
    echo "To analyze results:"
    echo "  go run cmd/analyze_memory_log.go $OUTPUT_DIR/stress_memory_leak_$TIMESTAMP.log --plot"
    echo "══════════════════════════════════════════════════════════════"
    
    # Run analysis
    echo ""
    echo "Running analysis..."
    go run cmd/analyze_memory_log.go "$OUTPUT_DIR/stress_memory_leak_$TIMESTAMP.log" --plot
else
    echo "Warning: Memory log file not found"
fi
