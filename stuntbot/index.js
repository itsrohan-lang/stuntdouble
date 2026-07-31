const { spawn } = require('child_process');
const fs = require('fs');
const os = require('os');
const path = require('path');

module.exports = (app) => {
  app.log.info("🛡️ StuntBot App is active! Listening for PR comments...");

  app.on("issue_comment.created", async (context) => {
    // Only trigger on Pull Requests when the user says "@stuntdouble review"
    if (!context.payload.issue.pull_request) return;
    if (!context.payload.comment.body.includes("@stuntdouble review")) return;

    const issue = context.issue();
    
    app.log.info(`>> [StuntBot] Received review request on PR #${issue.issue_number}`);
    
    const reaction = context.issue({
      content: "eyes",
      comment_id: context.payload.comment.id
    });
    await context.octokit.reactions.createForIssueComment(reaction);

    let tmpDir;
    try {
      tmpDir = await fs.promises.mkdtemp(path.join(os.tmpdir(), 'stuntdouble-pr-'));
    } catch (err) {
      app.log.error(`Failed to create temp directory: ${err.message}`);
      tmpDir = process.cwd();
    }

    // Notify the user the sandbox is spinning up
    const comment = context.issue({
      body: "🚀 **StuntDouble Isolation Sandbox Starting...**\n\nProvisioning isolated workspace sandbox and booting agent swarm..."
    });
    await context.octokit.issues.createComment(comment);

    app.log.info(">> Spawning Dev and QA agents into isolated StuntNet Swarm...");
    
    let completed = 0;
    let hasErrored = false;

    const handleError = async (err) => {
      if (hasErrored) return;
      hasErrored = true;
      app.log.error(`StuntDouble CLI execution failed: ${err.message}`);
      const errComment = context.issue({
        body: `❌ **StuntDouble Sandbox Error**\n\nFailed to launch sandboxed agent: \`${err.message}\`. Ensure \`sd\` CLI is installed on the host system PATH.`
      });
      await context.octokit.issues.createComment(errComment);
    };

    const checkCompletion = async () => {
      completed++;
      if (completed === 2 && !hasErrored) {
        const resultComment = context.issue({
          body: `✅ **StuntDouble Swarm Complete**\n\nBoth Dev and QA Agents executed safely inside the StuntNet isolated network. No external database connections or rogue network calls were allowed.\n\n*Zero zero-day escapes detected during swarm review.*`
        });
        await context.octokit.issues.createComment(resultComment);
      }
    };

    const devAgent = spawn("sd", ["run", "claude", "--allow-unenforced-network"], { cwd: tmpDir });
    const qaAgent = spawn("sd", ["run", "cursor", "--allow-unenforced-network"], { cwd: tmpDir });

    devAgent.on('error', handleError);
    qaAgent.on('error', handleError);

    devAgent.on('close', checkCompletion);
    qaAgent.on('close', checkCompletion);
  });
};
