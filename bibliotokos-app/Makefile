.PHONY: apple apple-setup apple-clean apple-verify help

APP_NAME = bibliotokos-app
BIN_DIR = bin
DMG_NAME = Bibliotokos-Installer.dmg

help:
	@echo "Bibliotokos Apple Distribution Makefile"
	@echo ""
	@echo "Before running 'make apple', you need to:"
	@echo "  1. Get an Apple Developer ID certificate (see README)"
	@echo "  2. Run 'make apple-setup' once to configure signing"
	@echo ""
	@echo "Available targets:"
	@echo "  apple-setup    - Interactive setup for code signing and notarization"
	@echo "  apple          - Build, sign, notarize, and create DMG"
	@echo "  apple-verify   - Verify the signed and notarized app"
	@echo "  apple-clean    - Clean build artifacts"

apple-setup:
	@echo "==> Setting up Apple code signing..."
	@echo "This will interactively configure:"
	@echo "  - Signing identity"
	@echo "  - Notarization credentials"
	@echo "  - Entitlements"
	@echo ""
	@echo "You will need:"
	@echo "  - Your Apple ID email"
	@echo "  - Your Team ID"
	@echo "  - An app-specific password from https://appleid.apple.com"
	@echo ""
	wails3 setup signing -platform darwin

apple:
	@echo "==> Building Bibliotokos for Apple Distribution..."
	@echo ""
	@echo "Step 1: Building universal binary (arm64 + amd64)..."
	wails3 task darwin:build:universal
	@echo ""
	@echo "Step 2: Creating and packaging .app bundle..."
	wails3 task darwin:package:universal
	@echo ""
	@echo "Step 3: Signing and notarizing with Apple..."
	wails3 task darwin:sign:notarize
	@echo ""
	@echo "Step 4: Creating DMG installer..."
	@$(MAKE) create-dmg
	@echo ""
	@echo "==> ✓ Build complete!"
	@echo "    Signed app: $(BIN_DIR)/$(APP_NAME).app"
	@echo "    Installer:  $(BIN_DIR)/$(DMG_NAME)"

create-dmg:
	@echo "Creating DMG installer..."
	@rm -f "$(BIN_DIR)/$(DMG_NAME)"
	@rm -rf "$(BIN_DIR)/dmg-temp"
	@mkdir -p "$(BIN_DIR)/dmg-temp"
	@cp -R "$(BIN_DIR)/$(APP_NAME).app" "$(BIN_DIR)/dmg-temp/Bibliotokos.app"
	@ln -s /Applications "$(BIN_DIR)/dmg-temp/Applications"
	@hdiutil create -volname "Bibliotokos" \
		-srcfolder "$(BIN_DIR)/dmg-temp" \
		-ov -format UDZO \
		"$(BIN_DIR)/$(DMG_NAME)"
	@rm -rf "$(BIN_DIR)/dmg-temp"
	@echo "DMG created with Applications symlink"

apple-verify:
	@echo "==> Verifying signature and notarization..."
	@echo ""
	@echo "Codesign verification:"
	@codesign -vvv --deep --strict "$(BIN_DIR)/$(APP_NAME).app"
	@echo ""
	@echo "Spctl assessment:"
	@spctl -a -vvv -t install "$(BIN_DIR)/$(APP_NAME).app"
	@echo ""
	@echo "Notarization status:"
	@stapler validate "$(BIN_DIR)/$(APP_NAME).app"

apple-clean:
	@echo "==> Cleaning build artifacts..."
	@rm -rf "$(BIN_DIR)/$(APP_NAME).app"
	@rm -rf "$(BIN_DIR)/Bibliotokos.app"
	@rm -f "$(BIN_DIR)/$(APP_NAME)"
	@rm -f "$(BIN_DIR)/$(APP_NAME)-amd64"
	@rm -f "$(BIN_DIR)/$(APP_NAME)-arm64"
	@rm -f "$(BIN_DIR)/$(DMG_NAME)"
	@echo "==> Clean complete!"
