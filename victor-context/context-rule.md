# Context File Usage Rules

## When Starting Work on a Feature

1. Load the appropriate context file FIRST
   - Core context: Always load all files in context/core/
   - Feature context: Load context/features/{feature-name}/context.md

2. Refer to context file for:
   - Architecture understanding
   - Code patterns and conventions
   - Relevant file locations
   - Common gotchas

3. Update context file when you learn something important:
   - Add to "WORKING NOTES" section immediately
   - Promote to main sections during daily review

## Context File Priority

Priority order when generating responses:
1. Context files (highest priority)
2. Actual codebase files
3. External documentation
4. General knowledge

## Keep Context Fresh

- Review context file at start of each day
- Update "LAST UPDATED" when making changes
- If context contradicts actual code, UPDATE CONTEXT (code is source of truth)