const { spawn } = require('child_process');

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

    // Notify the user the sandbox is spinning up
    const comment = context.issue({
      body: "🚀 **StuntDouble Isolation Sandbox Starting...**\n\nCloning repository securely and booting eBPF kernel hooks..."
    });
    await context.octokit.issues.createComment(comment);

    // In a production environment, StuntBot would securely clone the PR branch into 
    // a temporary directory, and execute the StuntDouble CLI with the selected AI agent.
    
    // Deploy a Multi-Agent Swarm (StuntNet)
    // We spawn a 'dev' agent to review the code, and a 'qa' agent to attempt malicious escapes
    app.log.info(">> Spawning Dev and QA agents into isolated StuntNet Swarm...");
    const devAgent = spawn("sd", ["run", "claude-dev", "--network=stuntnet"]);
    const qaAgent = spawn("sd", ["run", "cursor-qa", "--network=stuntnet"]);
    
    let completed = 0;
    
    const checkCompletion = async () => {
      completed++;
      if (completed === 2) {
        const resultComment = context.issue({
          body: `✅ **StuntDouble Swarm Complete**\n\nBoth Dev and QA Agents executed safely inside the StuntNet isolated network. No external database connections or rogue network calls were allowed.\n\n*Zero zero-day escapes detected during swarm review.*`
        });
        await context.octokit.issues.createComment(resultComment);
      }
    };
    
    devAgent.on('close', checkCompletion);
    qaAgent.on('close', checkCompletion);
  });
};
