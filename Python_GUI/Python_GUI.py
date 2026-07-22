import tkinter as tk
from tkinter import filedialog, ttk, messagebox
import subprocess
import os
import sys
# If you came here from probe.go, Look down at line 25 :)
# Base directory of the script or compiled executable
if getattr(sys, "frozen", False):
    APP_DIR = os.path.dirname(sys.executable)
else:
    APP_DIR = os.path.dirname(os.path.abspath(__file__))

# If running from dist/, go back to the project root.
# Otherwise stay in the script directory.

if os.path.basename(APP_DIR).lower() == "dist":
    BASE_DIR = os.path.dirname(APP_DIR)
else:
    BASE_DIR = APP_DIR


BIN_DIR = os.path.join(BASE_DIR, "bin")
EXE = ".exe" if os.name == "nt" else ""
RENDERERS = {
#Add the new Builds down below
    "Normal": os.path.join(BIN_DIR, f"maprenderer{EXE}"),
    "No Trees": os.path.join(BIN_DIR, "NOTREES", f"maprendererTREEGONE{EXE}"),
    "No Flowers": os.path.join(BIN_DIR, "NOFLOWER", f"maprendererNOFLOWER{EXE}"),
    "No Greenery": os.path.join(BIN_DIR, "NOGREENERY", f"maprendererNOBOTH{EXE}"),
}
def browse_folder():
    folder = filedialog.askdirectory()
    if folder:
        world_var.set(folder)

def run_render():
    renderer_name = renderer_var.get()
    exe = RENDERERS.get(renderer_name)

    if not exe or not os.path.exists(exe):
        messagebox.showerror("Error", "Renderer executable not found.")
        return

    world = world_var.get().strip()
    if not world:
        messagebox.showerror("Error", "Select a world folder.")
        return

    from_pos = f"{fx.get()},{fy.get()},{fz.get()}"
    to_pos = f"{tx.get()},{ty.get()},{tz.get()}"
    map_type = type_var.get()
    view = view_var.get()
    output = output_var.get().strip()

    PICTURES_DIR = os.path.join(BASE_DIR, "Pictures")
    os.makedirs(PICTURES_DIR, exist_ok=True)


    if not os.path.isabs(output):
        output = os.path.join(PICTURES_DIR, output)
    cmd = [
        exe,
        "-world", world,
        "-from", from_pos,
        "-to", to_pos,
        "-type", map_type,
        "-output", output
    ]

    # Add view only for isometric mode
    if map_type == "isometric":
        cmd.extend(["-view", view])

    status_var.set("Rendering...")
    root.update()

    try:
        result = subprocess.run(cmd, capture_output=True, text=True)

        if result.returncode == 0:
            status_var.set("Done.")
            messagebox.showinfo("Success", f"Render complete:\n{output}")
        else:
            status_var.set("Failed.")
            messagebox.showerror("Error", result.stderr or result.stdout)

    except Exception as e:
        status_var.set("Failed.")
        messagebox.showerror("Error", str(e))

root = tk.Tk()
root.title("MapRenderer GUI")
root.geometry("430x410")
root.resizable(False, False)

world_var = tk.StringVar()
renderer_var = tk.StringVar(value="Normal")
type_var = tk.StringVar(value="isometric")
view_var = tk.StringVar(value="nw")
output_var = tk.StringVar(value="map.png")
status_var = tk.StringVar(value="Ready")

# World Folder
tk.Label(root, text="World Folder").pack(anchor="w", padx=10, pady=(10,0))
frame1 = tk.Frame(root)
frame1.pack(fill="x", padx=10)

tk.Entry(frame1, textvariable=world_var, width=42).pack(side="left")
tk.Button(frame1, text="Browse", command=browse_folder).pack(side="left", padx=5)

# Renderer
tk.Label(root, text="Renderer").pack(anchor="w", padx=10, pady=(10,0))
ttk.Combobox(
    root,
    textvariable=renderer_var,
    values=list(RENDERERS.keys()),
    state="readonly"
).pack(fill="x", padx=10)

# Coordinates
tk.Label(root, text="From (X Y Z)(Desending values -ve)").pack(anchor="w", padx=10, pady=(10,0))
frame2 = tk.Frame(root)
frame2.pack(padx=10)

fx = tk.Entry(frame2, width=8); fx.insert(0, "9050"); fx.pack(side="left")
fy = tk.Entry(frame2, width=8); fy.insert(0, "150"); fy.pack(side="left", padx=5)
fz = tk.Entry(frame2, width=8); fz.insert(0, "27785"); fz.pack(side="left")

tk.Label(root, text="To (X Y Z)(Assending values +ve)").pack(anchor="w", padx=10, pady=(10,0))
frame3 = tk.Frame(root)
frame3.pack(padx=10)

tx = tk.Entry(frame3, width=8); tx.insert(0, "9250"); tx.pack(side="left")
ty = tk.Entry(frame3, width=8); ty.insert(0, "280"); ty.pack(side="left", padx=5)
tz = tk.Entry(frame3, width=8); tz.insert(0, "28130"); tz.pack(side="left")

# Map Type
tk.Label(root, text="Map Type").pack(anchor="w", padx=10, pady=(10,0))
ttk.Combobox(
    root,
    textvariable=type_var,
    values=["map", "isometric"],
    state="readonly"
).pack(fill="x", padx=10)

# View Direction
tk.Label(root, text="Isometric View").pack(anchor="w", padx=10, pady=(10,0))
ttk.Combobox(
    root,
    textvariable=view_var,
  values=["ne", "nw", "se", "sw"],
    state="readonly"
).pack(fill="x", padx=10)

# Output File
tk.Label(root, text="Output File").pack(anchor="w", padx=10, pady=(10,0))
tk.Entry(root, textvariable=output_var).pack(fill="x", padx=10)

# Render Button
tk.Button(root, text="Render", command=run_render, height=2).pack(fill="x", padx=10, pady=12)

# Status
tk.Label(root, textvariable=status_var, fg="blue").pack()

root.bind("<Return>", lambda event: run_render())
root.mainloop()
