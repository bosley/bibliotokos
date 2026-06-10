You are helping design the UI/UX for a lightweight personal Bible study desktop app built with Wails v3 and Svelte. The app is a focused, simpler alternative to Logos Bible Software.

### Core Purpose
The primary goals of the app are:
- Reading Scripture
- Comparing multiple Bible translations side-by-side
- Taking and managing personal notes
- Tracking theological concepts and connections (future)

### App Launch Behavior
- On startup, the app should apply the user’s saved theme preference (light/dark) from the system or user settings.

### Top Navigation Bar
Create a clean, minimalist top navigation bar with the following elements (from left to right):

- Two main icon buttons (SVG icons only):
  - **Scripture** icon
  - **Notes** icon
- Hovering over either icon should reveal its label as a tooltip.
- A search input field in the center.
- A theme toggle button (sun/moon icon) on the far right that switches between light and dark mode and persists the choice.

### Scripture View (Main Workspace)

Clicking the **Scripture** icon should open the main Bible reading workspace.

**Translation Selection:**
- The user should be able to select **one or more** Bible translations they have available.
- This selection can happen via a dropdown, segmented control, or similar UI element near the top of the workspace.

**Display Behavior:**

- **Single Translation Selected:**
  - The main content area should display the Scripture content for the user’s current search/range (e.g., typing “Gen” or “John 3” should load the full book or chapter).
  - Use a clean, readable verse-by-verse or paragraph layout.

- **Multiple Translations Selected:**
  - The main area should split into **side-by-side resizable panes** (similar to split views in VS Code or other code editors).
  - Each pane shows the **same search range** in a different translation.
  - Panes should be scrollable independently but stay synchronized in terms of the current book/chapter when possible.
  - Users should be able to add or remove panes dynamically.

The goal is to make cross-comparison of translations fast and natural.

### Notes Feature

Clicking the **Notes** icon should open a **separate window** (using Wails v3’s multi-window support).

**Notes Window Layout:**
- Left sidebar: A vertical list of all existing notes (showing note title/name).
- Main area: When a note is selected from the list, it opens in an editor view within the same window.
- The notes window should feel like a dedicated, focused note-taking space.

**Future Requirement (Note for later):**
- Notes should eventually be linkable to specific verses or passages in the Scripture view.

### Design Principles
- Keep the interface clean, calm, and highly readable.
- Prioritize clarity and low cognitive load.
- Use generous spacing and good typography for Scripture text.
- The UI should feel focused and purposeful rather than bloated.
