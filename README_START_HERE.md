# START HERE 👋

## Ready to Connect Mission Control to HubSpot?

I've set everything up for you! Here's what to do:

---

## 📋 **Quick Checklist**

### ✅ Already Done For You:
- ✅ Node.js installed (v25.2.1)
- ✅ HubSpot MCP server installed
- ✅ mission-control binary built and ready
- ✅ Setup scripts created

### ⏳ **You Need to Do** (10 minutes):
1. Create HubSpot OAuth app → Get Client ID
2. Run setup script → Enter Client ID
3. Start MCP server
4. Authenticate
5. Done! 🎉

---

## 🚀 The Easiest Way (Copy & Paste These Commands)

### **Command 1: Open HubSpot to Create OAuth App**

Open this URL in your browser:
```
https://developers.hubspot.com/
```

Then:
- Click **"Apps"** → **"Create app"** → **"Public app"**
- Name: `Mission Control CLI`
- Redirect URL: `http://127.0.0.1:8400/oauth/callback`
- Scopes: Select all 6 CRM scopes (contacts, companies, deals - read & write)
- Save and copy your **Client ID**

### **Command 2: Run Setup**

```bash
cd /Users/ericolden/workspace_launch/claude/mission-control
./quick-start.sh
```

Paste your Client ID when prompted.

### **Command 3: Start MCP Server (Keep Running)**

Open a new terminal and run:
```bash
cd /Users/ericolden/workspace_launch/claude/mission-control
source .env
PATH="/opt/homebrew/bin:$PATH" npx @hubspot/mcp-server
```

### **Command 4: Authenticate**

Back in your original terminal:
```bash
source .env
./mission-control auth login
```

Your browser will open → approve the app → done!

### **Command 5: Test It**

```bash
./mission-control auth status
./mission-control tools list
```

---

## 📚 **Detailed Guides Available**

- **NEXT_STEPS.md** - Detailed step-by-step walkthrough
- **SETUP_GUIDE.md** - Complete reference documentation
- **README.md** - Full mission-control documentation

---

## 🆘 **Need Help?**

Everything should work smoothly! If you hit any issues:

1. Make sure the MCP server is running (Command 3)
2. Make sure you ran `source .env` in your terminal
3. Check NEXT_STEPS.md for troubleshooting tips

---

## 🎯 **What You'll Be Able to Do**

Once set up, you can:

```bash
# Search for contacts
./mission-control hubspot contacts search --email john@example.com

# Create deals
./mission-control hubspot deals create --name "Big Deal" --amount 100000

# List all available HubSpot tools
./mission-control tools list

# Call any MCP tool directly
./mission-control tools call --name search_contacts --input '{"query":"test","limit":10}'
```

---

**Ready? Let's do this!** 🚀

Start with Command 1 above (create the HubSpot OAuth app).
