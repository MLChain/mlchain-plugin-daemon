class Mlchain < Formula
  desc "mlchain"
  homepage "https://github.com/mlchain/mlchain-plugin-daemon"
  version "0.0.1-beta.21"

  if OS.mac?
    if Hardware::CPU.intel?
      url "https://github.com/mlchain/mlchain-plugin-daemon/releases/download/0.0.1-beta.21/mlchain-plugin-darwin-amd64"
    elsif Hardware::CPU.arm?
      url "https://github.com/mlchain/mlchain-plugin-daemon/releases/download/0.0.1-beta.21/mlchain-plugin-darwin-arm64"
    end
  elsif OS.linux?
    if Hardware::CPU.intel?
      url "https://github.com/mlchain/mlchain-plugin-daemon/releases/download/0.0.1-beta.21/mlchain-plugin-linux-amd64"
    elsif Hardware::CPU.arm?
      url "https://github.com/mlchain/mlchain-plugin-dmlchain aemon/releases/download/0.0.1-beta.21/mlchain-plugin-linux-arm64"
    end
  elsif OS.windows?
    url "https://github.com/mlchain/mlchain-plugin-daemon/releases/download/0.0.1-beta.21/mlchain-plugin-windows-amd64"
  end

  def install
    bin.install "mlchain-plugin-darwin-#{Hardware::CPU.arch}" => "mlchain"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/mlchain --version")
  end
end
