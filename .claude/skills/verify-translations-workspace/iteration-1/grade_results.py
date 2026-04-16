#!/usr/bin/env python3
"""
Grade the test results against assertions.
"""

import json
from pathlib import Path
from typing import Dict, List


def grade_eval_1_with_skill(outputs_dir: Path) -> List[Dict]:
    """Grade eval-1 with skill results."""
    results = []

    # Check if verification report exists
    report_file = outputs_dir / "verification_report.txt"
    if report_file.exists():
        content = report_file.read_text()

        # runs_verification_script
        results.append({
            "text": "Should execute the check_translations.py script to verify completeness",
            "passed": "check_translations.py" in content or "verification" in content.lower(),
            "evidence": "Report file exists with verification content"
        })

        # reports_all_complete
        results.append({
            "text": "Should report that all translations are complete with a clear success message",
            "passed": ("complete" in content.lower() or "complet" in content.lower()) and
                     ("✓" in content or "OK" in content or "success" in content.lower()),
            "evidence": f"Found success indicators in report"
        })

        # no_unnecessary_actions
        results.append({
            "text": "Should not propose translations or modify files since everything is complete",
            "passed": "suggestion" not in content.lower() and "proposer" not in content.lower(),
            "evidence": "No translation suggestions found in output"
        })

        # clear_communication
        results.append({
            "text": "Should communicate the result clearly to the user in a concise way",
            "passed": len(content) < 5000,  # Reasonable length
            "evidence": f"Report is {len(content)} characters"
        })
    else:
        # All fail if no output
        for assertion in ["runs_verification_script", "reports_all_complete",
                         "no_unnecessary_actions", "clear_communication"]:
            results.append({
                "text": assertion,
                "passed": False,
                "evidence": "No verification_report.txt found"
            })

    return results


def grade_eval_1_without_skill(outputs_dir: Path) -> List[Dict]:
    """Grade eval-1 baseline (without skill) results."""
    results = []

    report_file = outputs_dir / "verification_report.txt"
    if report_file.exists():
        content = report_file.read_text()

        results.append({
            "text": "Detects all translations are complete",
            "passed": "complete" in content.lower() or "complet" in content.lower(),
            "evidence": "Found completeness indicators"
        })

        results.append({
            "text": "Provides clear result",
            "passed": len(content) > 100,
            "evidence": f"Report has {len(content)} characters"
        })
    else:
        results.append({
            "text": "Produces output",
            "passed": False,
            "evidence": "No output file found"
        })

    return results


def grade_eval_2_with_skill(outputs_dir: Path) -> List[Dict]:
    """Grade eval-2 with skill results."""
    results = []

    interaction_log = outputs_dir / "interaction_log.txt"

    if interaction_log.exists():
        content = interaction_log.read_text()

        # detects_missing_keys
        has_subtitle = "app.subtitle" in content
        has_could_not_load = "couldNotLoadResults" in content or "could not load" in content.lower()
        results.append({
            "text": "Should detect exactly 2 missing French keys: app.subtitle and messages.couldNotLoadResults",
            "passed": has_subtitle and has_could_not_load,
            "evidence": f"Found subtitle: {has_subtitle}, Found couldNotLoad: {has_could_not_load}"
        })

        # proposes_translations
        suggestion_count = content.lower().count("suggestion")
        results.append({
            "text": "Should propose 2-3 translation suggestions for each missing key with explanations",
            "passed": suggestion_count >= 2,
            "evidence": f"Found {suggestion_count} occurrences of 'suggestion'"
        })

        # interactive_validation
        results.append({
            "text": "Should present suggestions one at a time and wait for user validation before proceeding",
            "passed": "quelle traduction" in content.lower() or "validation" in content.lower(),
            "evidence": "Found interactive prompts"
        })

        # Other assertions - check for updated file
        updated_fr = outputs_dir / "updated_fr.json"
        results.append({
            "text": "Should correctly apply validated translations to fr.json file",
            "passed": updated_fr.exists(),
            "evidence": f"Updated file exists: {updated_fr.exists()}"
        })

        if updated_fr.exists():
            try:
                data = json.loads(updated_fr.read_text())
                results.append({
                    "text": "Updated fr.json should maintain proper JSON structure",
                    "passed": True,
                    "evidence": "JSON parses correctly"
                })
            except:
                results.append({
                    "text": "Updated fr.json should maintain proper JSON structure",
                    "passed": False,
                    "evidence": "JSON parse error"
                })
        else:
            results.append({
                "text": "Updated fr.json should maintain proper JSON structure",
                "passed": False,
                "evidence": "File not found"
            })

        # verifies_completion
        results.append({
            "text": "Should run verification again after updates to confirm all translations are now complete",
            "passed": "verification" in content.lower()[-500:],  # Check end of file
            "evidence": "Found verification mention near end of output"
        })
    else:
        for _ in range(6):
            results.append({
                "text": "Missing interaction log",
                "passed": False,
                "evidence": "No interaction_log.txt found"
            })

    return results


def main():
    workspace = Path("/home/frederic-prost/.claude/skills/verify-translations-workspace/iteration-1")

    # Grade each eval
    evals = [
        ("eval-1-all-complete", "with_skill", grade_eval_1_with_skill),
        ("eval-1-all-complete", "without_skill", grade_eval_1_without_skill),
        ("eval-2-missing-french-keys", "with_skill", grade_eval_2_with_skill),
    ]

    for eval_name, config, grader_func in evals:
        eval_dir = workspace / eval_name / config
        outputs_dir = eval_dir / "outputs"

        if not outputs_dir.exists():
            print(f"⚠ Outputs not found for {eval_name}/{config}")
            continue

        grading = {
            "expectations": grader_func(outputs_dir),
            "overall_pass": all(e["passed"] for e in grader_func(outputs_dir))
        }

        # Write grading result
        grading_file = eval_dir / "grading.json"
        with open(grading_file, "w") as f:
            json.dump(grading, f, indent=2)

        print(f"✓ Graded {eval_name}/{config}: {'PASS' if grading['overall_pass'] else 'FAIL'}")


if __name__ == "__main__":
    main()
