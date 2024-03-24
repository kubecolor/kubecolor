class Kubecolor < Formula
  desc "Colorize your kubectl output"
  homepage "https://kubecolor.github.io/"
  url "https://github.com/kubecolor/kubecolor/archive/refs/tags/v0.2.2.tar.gz"
  sha256 "ba0894a8e26fefff47a0691529964303bdd8fdc2d7ce74e7d241cb5a2f2ade50"
  license "MIT"

  depends_on "go" => :build
  depends_on "kubernetes-cli" => :test

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w")
  end

  test do
    system bin/"kubecolor", "-h"
    assert_equal 0, $CHILD_STATUS.exitstatus
  end
end
