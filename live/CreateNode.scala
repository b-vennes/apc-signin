//> using scala lts
//> using dep "org.typelevel::mouse:1.3.2"
//> using dep "org.typelevel::cats-effect:3.6.3"
//> using dep "com.lihaoyi::os-lib:0.11.5"
//> using dep "com.monovore::decline-effect:2.5.0"

package main

import cats.Show
import cats.data.NonEmptyList
import cats.implicits._
import cats.effect.{IO, IOApp}
import cats.effect.ExitCode
import com.monovore.decline.effect.CommandIOApp
import com.monovore.decline.Opts
import com.monovore.decline.Argument
import os.Path
import cats.data.ValidatedNel
import cats.data.Validated

object Display:
  def oneLine(textBlock: String): String =
    textBlock
      .stripMargin
      .split(System.lineSeparator())
      .mkString(" ")

object DigitalOcean:
  object Droplet:
    opaque type Image = String
    object Image:
      val ubuntu: Image = "ubuntu-25-04-x64"

      given Show[Image] = s => s

    opaque type Size = String
    object Size:
      val level2: Size = "s-1vcpu-2gb"

      given Show[Size] = s => s

    opaque type Region = String
    object Region:
      val sfo: Region = "sfo3"

      given Show[Region] = r => r

  def makeDroplet(
    image: Droplet.Image,
    size: Droplet.Size,
    region: Droplet.Region,
    projectID: String,
    vpc: String,
    sshKeyID: String,
    name: String
  ): IO[Unit] = IO.blocking:
    os.proc(
      "doctl",
      "compute",
      "droplet",
      "create",
      name,
      "--image",
      image.show,
      "--size",
      size.show,
      "--region",
      region.show,
      "--vpc-uuid",
      vpc.show,
      "--ssh-keys",
      sshKeyID,
      "--project-id",
      projectID,
    ).call(
      cwd = os.pwd,
      stdin = os.Inherit,
      stdout = os.Inherit,
      stderr = os.Inherit
    )
  .flatMap: result =>
    IO.println(s"Result of make droplet is '${result.exitCode}'.")

object CreateNode extends CommandIOApp(
  "scala-cli live/CreateNode.scala",
  "Create a new Droplet node in Digital Ocean.",
  helpFlag = true,
  "0.0.1"
):
  given Argument[os.Path] = new Argument[os.Path]:
    def defaultMetavar: String = "file path"
    def read(string: String): ValidatedNel[String, Path] =
      Validated
        .catchNonFatal(os.Path(string))
        .orElse(Validated.catchNonFatal(os.RelPath(string)).map(os.pwd / _))
        .leftMap(error => NonEmptyList.of("Invalid path."))

  def main: Opts[IO[ExitCode]] =
    (
      Opts.env[String](
        "VPC_ID",
        Display.oneLine(
          """The ID of the VPC to create the droplet in.
            |Find with 'doctl vpcs list --format REGION,ID'."""
        )
      ),
      Opts.env[String](
        "SSH_KEY",
        Display.oneLine(
          """An SSH key ID to initialize the compute with.
            |Find with 'doctl compute ssh-key list'."""
        )
      ),
      Opts.env[String](
        "COMPUTE_NAME",
        Display.oneLine(
          """The name of the new compute instance."""
        )
      ),
      Opts.env[String](
        "PROJECT_ID",
        Display.oneLine(
          """The project ID this droplet will be placed in.
            |Find with 'doctl projects list --format Name,ID'."""
        )
      )
    )
      .mapN: (vpcID, sshKey, computeName, projectID) =>
        DigitalOcean.makeDroplet(
          image = DigitalOcean.Droplet.Image.ubuntu,
          size = DigitalOcean.Droplet.Size.level2,
          region = DigitalOcean.Droplet.Region.sfo,
          projectID = projectID,
          vpc = vpcID,
          sshKeyID = sshKey,
          name = computeName
        ) >>
        IO(ExitCode.Success)
