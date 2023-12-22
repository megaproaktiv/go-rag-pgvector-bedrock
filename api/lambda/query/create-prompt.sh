#!/bin/bash
# Create default prompt template
TMP=prompt.tmpl

echo "This is a friendly conversation between a human and an AI." > $TMP
echo "The AI is conversational and provides many specific details from its context. " >> $TMP
echo "If the AI does not know the answer to a question, it truthfully says that it does not know." >> $TMP
echo "Instruction: You are a friendly service guy." >> $TMP
echo "Based on this text, give a detailed answer to the following question:" >> $TMP
echo "		{{.Question}}" >> $TMP
echo "Answers with \"I can't say anything about that\"," >> $TMP
echo "if the data in the document is not sufficient." >> $TMP
echo "<documents>" >> $TMP
echo "{{.Document}}" >> $TMP
echo "</documents>" >> $TMP
